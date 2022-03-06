package apiserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	logger "github.com/webdevolegkuprianov/server_http_rest/app/apiserver/logger"
	"github.com/webdevolegkuprianov/server_http_rest/app/apiserver/model"
	"github.com/webdevolegkuprianov/server_http_rest/app/apiserver/store"
)

//errors
var (
	errIncorrectEmailOrPassword = errors.New("incorrect auth")
	errReg                      = errors.New("service registration error")
	errJwt                      = errors.New("token error")
	errFindUser                 = errors.New("user not found")
	errMssql                    = errors.New("mssql error")
)

//responses
var (
	respGazCrmWorkList = "data work_list recieved"
	respGazCrmLeadGet  = "data lead_get recieved"
	respGazCrmStatuses = "data statuses recieved"
	respBooking        = "data booking sent to gazcrm"
	respForm           = "data form sent to gazcrm"
	errPg              = "error postgres storing"
)

//server configure
type server struct {
	router *mux.Router
	store  store.Store
	config *model.Service
	client *http.Client
}

func newServer(store store.Store, config *model.Service, client *http.Client) *server {
	s := &server{
		router: mux.NewRouter(),
		store:  store,
		config: config,
		client: client,
	}
	s.configureRouter()
	return s
}

//write new token struct
func newToken(token string, exp time.Time) *model.Token_exp {
	return &model.Token_exp{
		Token: token,
		Exp:   exp,
	}
}

//write response struct
func newResponse(status string, response string) *model.Response {
	return &model.Response{
		Status:   status,
		Response: response,
	}
}

//write response struct booking
func newResponseBooking(statusms string, responsems string, statusgcrm string, responsegcrm string) *model.ResponseBooking {
	return &model.ResponseBooking{
		StatusMs:       statusms,
		ResponseMs:     responsems,
		StatusGazCrm:   statusgcrm,
		ResponseGazCrm: responsegcrm,
	}
}

//write http error
func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})

}

//write http response
func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	//open
	s.router.HandleFunc("/authentication", s.handleAuth()).Methods("POST")
	//private
	auth := s.router.PathPrefix("/auth").Subrouter()
	auth.Use(s.middleWare)
	//booking, forms submit
	auth.HandleFunc("/requestbooking", s.handleRequestBooking()).Methods("POST")
	auth.HandleFunc("/requestform", s.handleRequestForm()).Methods("POST")
	//gaz crm
	auth.HandleFunc("/requestleadget", s.handleRequestLeadGetGazCrm()).Methods("POST")
	auth.HandleFunc("/requestworklist", s.handleRequestWorkListGazCrm()).Methods("POST")
	auth.HandleFunc("/requeststatus", s.handleRequestStatusGazCrm()).Methods("POST")
	//stock
	auth.HandleFunc("/getdatastocks", s.handleGetDataStocks()).Methods("GET")
	//prices
	auth.HandleFunc("/getbasicmodelsprice", s.handleBasicModelsPrice()).Methods("GET")
	auth.HandleFunc("/getoptionsprice", s.handleOptionsPrice()).Methods("GET")
	auth.HandleFunc("/getgeneralprice", s.handleGeneralPrice()).Methods("GET")
	//sprav models
	auth.HandleFunc("/getsprav", s.handleSprav()).Methods("GET")
	//options
	auth.HandleFunc("/getoptionsdata", s.handleOptionsData()).Methods("GET")
	auth.HandleFunc("/getoptionsdatasprav", s.handleOptionsDataSprav()).Methods("GET")
	auth.HandleFunc("/getpacketsdata", s.handlePacketsData()).Methods("GET")
	//colors
	auth.HandleFunc("/getcolorsdata", s.handleColorsData()).Methods("GET")
}

//handle Auth
func (s *server) handleAuth() http.HandlerFunc {

	var req model.User1

	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, errReg)
			logger.ErrorLogger.Println(err)
			return
		}

		u, err := s.store.User().FindUser(req.Email, req.Password)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			logger.ErrorLogger.Println(err)
			return
		}

		token, datetime_exp, err := s.store.User().CreateToken(uint64(u.ID), s.config)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, errJwt)
			logger.ErrorLogger.Println(err)
			return
		}
		token_data := newToken(token, datetime_exp)
		s.respond(w, r, http.StatusOK, token_data)
		logger.InfoLogger.Println("token issued success")

	}

}

//Middleware
func (s *server) middleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//extract user_id
		user_id, err := s.store.User().ExtractTokenMetadata(r, s.config)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errJwt)
			logger.ErrorLogger.Println(err)
			return
		}

		if err := s.store.User().FindUserid(user_id.UserId); err != nil {
			s.error(w, r, http.StatusUnauthorized, errFindUser)
			logger.ErrorLogger.Println(err)
			return
		}

		next.ServeHTTP(w, r)

	})

}

//handle Client Data
func (s *server) handleRequestBooking() http.HandlerFunc {

	var errMs string

	return func(w http.ResponseWriter, r *http.Request) {

		req := model.DataBooking{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			logger.ErrorLogger.Println(err)
			return
		}

		resp, err := s.store.Data().QueryInsertMssql(req)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, errMssql)
			logger.ErrorLogger.Println(err)
			logger.ErrorLogger.Println(resp)
			return
		}

		if resp != "Обработка данных прошла успешно" {
			errMs = "Error"
			logger.ErrorLogger.Println(resp)
		} else {
			errMs = "Ok"
			logger.InfoLogger.Println("data booking stored in mssql")

			//respm, err := s.store.Data().CallMSMailing(req, s.config)
			//if err != nil {
			//ErrorLogger.Println(err)
			//ErrorLogger.Println(respm)
			//}
			//InfoLogger.Println("email=" + respm)

		}

		//request gazcrm api
		respg, err := s.store.Data().RequestGazCrmApiBooking(req, s.config)
		if err != nil {
			logger.ErrorLogger.Println(err)
		}
		if respg.Status != "OK" {
			logger.ErrorLogger.Println(respg)
			s.respond(w, r, http.StatusBadRequest, newResponseBooking(errMs, resp, "Error", respg.Message))
		} else {
			logger.InfoLogger.Println("gazcrm booking data transfer success")
			s.respond(w, r, http.StatusOK, newResponseBooking(errMs, resp, "Ok", respBooking))
		}

		//insert data in postgres
		if err := s.store.Data().QueryInsertBookingPostgres(req); err != nil {
			logger.ErrorLogger.Println(err)
		} else {
			logger.InfoLogger.Println("sites booking data stored")
		}

	}

}

//handle request forms
func (s *server) handleRequestForm() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		req := model.DataForms{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			logger.ErrorLogger.Println(err)
			return
		}

		//request gazcrm api
		respg, err := s.store.Data().RequestGazCrmApiForms(req, s.config)
		if err != nil {
			logger.ErrorLogger.Println(err)
		}
		if respg.Status != "OK" {
			logger.ErrorLogger.Println(respg)
			s.respond(w, r, http.StatusBadRequest, newResponse("Error", respg.Message))
		} else {
			logger.InfoLogger.Println("gazcrm form data transfer success")
			s.respond(w, r, http.StatusOK, newResponse("Ok", respForm))
		}

		//insert data in postgres
		if err := s.store.Data().QueryInsertFormsPostgres(req); err != nil {
			logger.ErrorLogger.Println(err)
		} else {
			logger.InfoLogger.Println("sites form data stored")
		}

	}

}

//gaz crm
//handle request lead get from gaz crm
func (s *server) handleRequestLeadGetGazCrm() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		req := model.DataLeadGet{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			logger.ErrorLogger.Println(err)
			return
		}

		//insert data in postgres
		if err := s.store.Data().QueryInsertLeadGetPostgres(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			logger.ErrorLogger.Println(err)
			s.respond(w, r, http.StatusBadRequest, newResponse("Error", errPg))
		} else {
			logger.InfoLogger.Println("gazcrm lead_get inserted in postgres")
			s.respond(w, r, http.StatusOK, newResponse("Ok", respGazCrmLeadGet))
		}

	}

}

//handle request work list from gaz crm
func (s *server) handleRequestWorkListGazCrm() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		req := model.DataWorkList{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			logger.ErrorLogger.Println(err)
			return
		}

		//insert data in postgres
		if err := s.store.Data().QueryInsertWorkListsPostgres(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			logger.ErrorLogger.Println(err)
			s.respond(w, r, http.StatusBadRequest, newResponse("Error", errPg))
		} else {
			logger.InfoLogger.Println("gazcrm work_list inserted in postgres")
			s.respond(w, r, http.StatusOK, newResponse("Ok", respGazCrmWorkList))
		}

	}

}

//handle request status from gaz crm
func (s *server) handleRequestStatusGazCrm() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		req := model.DataStatuses{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			logger.ErrorLogger.Println(err)
			return
		}

		//insert data in postgres
		if err := s.store.Data().QueryInsertStatusesPostgres(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			logger.ErrorLogger.Println(err)
			s.respond(w, r, http.StatusBadRequest, newResponse("Error", errPg))
		} else {
			logger.InfoLogger.Println("gazcrm statuses inserted in postgres")
			s.respond(w, r, http.StatusOK, newResponse("Ok", respGazCrmStatuses))
		}

	}

}

//handle request stocks
func (s *server) handleGetDataStocks() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		data, err := s.store.Data().QueryStocksMssql()

		if err != nil {
			s.error(w, r, http.StatusBadRequest, errMssql)
			logger.ErrorLogger.Println(err)
		}

		s.respond(w, r, http.StatusOK, data)
		logger.InfoLogger.Println("data stocks sent")

	}

}

//handle request basic model price
func (s *server) handleBasicModelsPrice() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		data, err := s.store.Data().QueryBasicModelsPriceMssql()

		if err != nil {
			s.error(w, r, http.StatusBadRequest, errMssql)
			logger.ErrorLogger.Println(err)
		}

		s.respond(w, r, http.StatusOK, data)
		logger.InfoLogger.Println("data price basic models sent")

	}

}

//handle request options price
func (s *server) handleOptionsPrice() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		data, err := s.store.Data().QueryOptionsPriceMssql()

		if err != nil {
			s.error(w, r, http.StatusBadRequest, errMssql)
			logger.ErrorLogger.Println(err)
		}

		s.respond(w, r, http.StatusOK, data)
		logger.InfoLogger.Println("data price options sent")

	}

}

//handle request general price
func (s *server) handleGeneralPrice() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		data, err := s.store.Data().QueryGeneralPriceMssql()

		if err != nil {
			s.error(w, r, http.StatusBadRequest, errMssql)
			logger.ErrorLogger.Println(err)
		}

		s.respond(w, r, http.StatusOK, data)
		logger.InfoLogger.Println("data price general sent")

	}

}

//handle request sprav
func (s *server) handleSprav() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		data, err := s.store.Data().QuerySprav()

		if err != nil {
			s.error(w, r, http.StatusBadRequest, errMssql)
			logger.ErrorLogger.Println(err)
		}

		s.respond(w, r, http.StatusOK, data)
		logger.InfoLogger.Println("data sprav sent")

	}

}

//handle request options data
func (s *server) handleOptionsData() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		data, err := s.store.Data().QueryOptionsData()

		if err != nil {
			s.error(w, r, http.StatusBadRequest, errMssql)
			logger.ErrorLogger.Println(err)
		}

		s.respond(w, r, http.StatusOK, data)
		logger.InfoLogger.Println("data options sent")

	}

}

//handle request options sprav data
func (s *server) handleOptionsDataSprav() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		data, err := s.store.Data().QueryOptionsDataSprav()

		if err != nil {
			s.error(w, r, http.StatusBadRequest, errMssql)
			logger.ErrorLogger.Println(err)
		}

		s.respond(w, r, http.StatusOK, data)
		logger.InfoLogger.Println("data options sprav sent")

	}

}

//handle request options packets data
func (s *server) handlePacketsData() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		data, err := s.store.Data().QueryPacketsData()

		if err != nil {
			s.error(w, r, http.StatusBadRequest, errMssql)
			logger.ErrorLogger.Println(err)
		}

		s.respond(w, r, http.StatusOK, data)
		logger.InfoLogger.Println("data packets sent")

	}

}

//handle request colors data
func (s *server) handleColorsData() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		data, err := s.store.Data().QueryColorsData()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, errMssql)
			logger.ErrorLogger.Println(err)
		}

		s.respond(w, r, http.StatusOK, data)
		logger.InfoLogger.Println("data colors sent")

	}

}

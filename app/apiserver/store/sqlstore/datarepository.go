package sqlstore

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	"github.com/webdevolegkuprianov/server_http_rest/app/apiserver/model"

	logger "github.com/webdevolegkuprianov/server_http_rest/app/apiserver/logger"
)

//Data repository
type DataRepository struct {
	store *Store
}

//query insert mssql
func (r *DataRepository) QueryInsertMssql(data model.DataBooking) (string, error) {

	//validation
	if err := data.ValidateDataBooking(); err != nil {
		logger.ErrorLogger.Println(err)
		return "", err
	}

	//request mssql
	var mssql_respond string

	_, err := r.store.dbMssql.Exec(r.store.config.Spec.Queryies.Booking,
		sql.Named("ИдентификаторОбращения", data.RequestId),
		sql.Named("Действие", data.ActionType),
		sql.Named("НомернойТовар", data.UniqModCode),
		sql.Named("ВИН", data.Vin),
		sql.Named("ЦенаСНДС", data.PriceWithNds),
		sql.Named("ТипКонтрагента", data.TypeClient),
		sql.Named("ИНН", data.Inn),
		sql.Named("КПП", data.Kpp),
		sql.Named("ОГРН", data.Ogrn),
		sql.Named("АдресЮридический", data.YurAddressCode),
		sql.Named("АдресДоставки", data.DeliveryAddressCode),
		sql.Named("Hid", data.Hid),
		sql.Named("Наименование", data.CompanyName),
		sql.Named("Фамилия", data.Surname),
		sql.Named("Имя", data.Name),
		sql.Named("ДатаРождения", data.DateOfBirth),
		sql.Named("Отчество", data.Patronymic),
		sql.Named("СерияПаспорта", data.PassportSer),
		sql.Named("НомерПаспорта", data.PassportNumber),
		sql.Named("СНИЛС", data.Snils),
		sql.Named("ЭлектроннаяПочта", data.Email),
		sql.Named("Телефоны", data.PhoneNumber),
		sql.Named("МоментОбращения", data.TimeRequest),
		sql.Named("НомерСчета", data.BillNumber),
		sql.Named("Ошибка", sql.Out{Dest: &mssql_respond}),
		sql.Named("ВыполнитьТестовыйВызов", data.TestMod),
	)
	if err != nil {
		return "", err
	}

	return mssql_respond, nil
}

//insert booking in postgres
func (r *DataRepository) QueryInsertBookingPostgres(data model.DataBooking) error {

	query := `
	insert into booking
	values($1, $2, $3, $4, $5, $6, $7, $8, $9,
		$10, $11, $12, $13, $14, $15, $16, $17, $18,
		$19, $20, $21, $22, $23, $24, $25, $26, $27,
		$28, $29, $30, $31, $32, $33, $34, $35, $36,
		$37, $38, $39)`

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	_, err = tx.Exec(ctx, query,
		data.RequestId,
		data.ActionType,
		data.UniqModCode,
		data.Modification,
		data.ModFamily,
		data.ModBodyType,
		data.ModEngine,
		data.ModBase,
		data.ModTuning,
		data.Vin,
		data.PriceWithNds,
		data.TypeClient,
		data.Inn,
		data.Kpp,
		data.Ogrn,
		data.YurAddressCode,
		data.DeliveryAddressCode,
		data.DeliveryAddress,
		data.Hid,
		data.CompanyName,
		data.RepresentativeName,
		data.RepresentativeSurname,
		data.Surname,
		data.Name,
		data.Patronymic,
		data.PassportSer,
		data.PassportNumber,
		data.Snils,
		data.DateOfBirth,
		data.Email,
		data.PhoneNumber,
		data.Comment,
		data.Consentmailing,
		data.TimeRequest,
		data.File,
		data.BillNumber,
		data.UrlMod,
		data.Clientid,
		data.Ymuid,
		data.TestMod,
	)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	return nil

}

//insert forms in postgres
func (r *DataRepository) QueryInsertFormsPostgres(data model.DataForms) error {

	query := `
	insert into forms
	values($1, $2, $3, $4, $5, $6, $7, $8, $9,
		$10, $11, $12, $13, $14, $15, $16, $17, $18,
		$19, $20, $21, $22, $23, $24, $25, $26, $27,
		$28, $29, $30, $31)`

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	_, err = tx.Exec(ctx, query,
		data.TimeRequest,
		data.RequestId,
		data.SubdivisionsId,
		data.SubdivisionsName,
		data.FormName,
		data.FormId,
		data.HostName,
		data.Division,
		data.Area,
		data.BrandName,
		data.CarModel,
		data.Clientid,
		data.MetricsType,
		data.СlientIP,
		data.TypeClient,
		data.CompanyName,
		data.Name,
		data.Email,
		data.PhoneNumber,
		data.Comment,
		data.Consentmailing,
		data.ActionType,
		data.Modification,
		data.ModFamily,
		data.ModBodyType,
		data.ModEngine,
		data.ModBase,
		data.ModTuning,
		data.Vin,
		data.PriceWithNds,
		data.UrlMod,
	)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	return nil

}

//insert lead_get in postgres
func (r *DataRepository) QueryInsertLeadGetPostgres(data model.DataLeadGet) error {

	valSlice := reflect.ValueOf(data).FieldByName("Data").Interface().([]model.DataLeadGet_Gazcrm)

	query := `
	insert into gazcrm_lead_get
	values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	_, err = tx.Exec(ctx, query,
		valSlice[0].TimeRequest,
		valSlice[1].EventName,
		valSlice[2].RequestId,
		valSlice[3].SubdivisionsId,
		valSlice[4].SubdivisionsName,
		valSlice[5].FormName,
		valSlice[6].HostName,
		valSlice[7].Division,
		valSlice[8].Area,
		valSlice[9].BrandName,
		valSlice[10].ClientID,
		valSlice[11].MetricsType,
	)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	return nil
}

//insert work lists in postgres
func (r *DataRepository) QueryInsertWorkListsPostgres(data model.DataWorkList) error {

	valSlice := reflect.ValueOf(data).FieldByName("Data").Interface().([]model.DataWorkList_Gazcrm)

	query := `
	insert into gazcrm_work_list
	values($1, $2, $3, $4)`

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	_, err = tx.Exec(ctx, query,
		valSlice[0].TimeRequest,
		valSlice[1].EventName,
		valSlice[2].GazcrmClientId,
		valSlice[3].GazCrmWorkListId,
	)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	return nil
}

//insert work lists in postgres
func (r *DataRepository) QueryInsertStatusesPostgres(data model.DataStatuses) error {

	valSlice := reflect.ValueOf(data).FieldByName("Data").Interface().([]model.DataStatuses_Gazcrm)

	query := `
	insert into gazcrm_statuses
	values($1, $2, $3, $4, $5, $6, $7)`

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	_, err = tx.Exec(ctx, query,
		valSlice[0].TimeRequest,
		valSlice[1].EventName,
		valSlice[2].RequestId,
		valSlice[3].GazcrmClientId,
		valSlice[4].GazCrmWorkListId,
		valSlice[5].ClientID,
		valSlice[6].MetricsType,
	)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	return nil
}

//query stocks mssql
func (r *DataRepository) QueryStocksMssql() ([]model.DataStocks, error) {

	rows, err := r.store.dbMssql.Query(r.store.config.Spec.Queryies.Stocks)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	defer rows.Close()

	results := []model.DataStocks{} // creating empty slice

	for rows.Next() {

		data := &model.DataStocks{} // creating new struct for every row

		err := rows.Scan(
			&data.VIN,
			&data.Площадка,
			&data.Наименование_номенклатуры,
			&data.Номер_согласно_КД,
			&data.Дивизион,
			&data.Доработчик_Подрядчик,
			&data.Test_truck,
			&data.Телематика,
			&data.Номер_шасси,
			&data.Номер_двигателя,
			&data.Грузоподъемность_кг,
			&data.Цвет,
			&data.Вариант_сборки,
			&data.Расшифровка_варианта_сборки,
			&data.Вариант_сборки_свернутый,
			&data.Год_VIN,
			&data.Дата_сборки,
			&data.Справочная_стоимость_по_прайсу,
			&data.Дата_отгрузки,
			&data.Дата_прихода,
			&data.Страна,
			&data.Контрагент_получателя,
			&data.Стоянка,
			&data.Город_стоянки,
			&data.Площадка_получателя_Ид,
			&data.Контрагент_получателя_Ид,
			&data.Город_стоянки_Ид,
			&data.Номер_заявки,
			&data.Для_доработки,
			&data.Номерной_товар,
		)
		if err != nil {
			logger.ErrorLogger.Println(err)
			return nil, err
		}
		results = append(results, *data)
	}

	return results, nil

}

//query mssql price basic models
func (r *DataRepository) QueryBasicModelsPriceMssql() ([]model.DataBasicModelsPrice, error) {

	rows, err := r.store.dbMssql.Query(r.store.config.Spec.Queryies.BasicModelsPrice)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	defer rows.Close()

	results := []model.DataBasicModelsPrice{} // creating empty slice

	for rows.Next() {

		data := &model.DataBasicModelsPrice{} // creating new struct for every row

		err := rows.Scan(
			&data.Товар,
			&data.Цена,
			&data.СтавкаНДС,
			&data.НДС,
			&data.НачалоДействия,
		)
		if err != nil {
			logger.ErrorLogger.Println(err)
			return nil, err
		}
		results = append(results, *data)
	}

	return results, nil

}

//query mssql options price
func (r *DataRepository) QueryOptionsPriceMssql() ([]model.DataOptionsPrice, error) {

	rows, err := r.store.dbMssql.Query(r.store.config.Spec.Queryies.OptionsPrice)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	defer rows.Close()

	results := []model.DataOptionsPrice{} // creating empty slice

	for rows.Next() {

		data := &model.DataOptionsPrice{} // creating new struct for every row

		err := rows.Scan(
			&data.ЕНСП_Модификация_Ид,
			&data.Товар,
			&data.ЗначениеОпции,
			&data.ОбозначениеОпции,
			&data.Цена,
			&data.СтавкаНДС_Ид,
			&data.НДС,
			&data.НачалоДействия,
			&data.СоставПакета,
		)
		if err != nil {
			logger.ErrorLogger.Println(err)
			return nil, err
		}
		results = append(results, *data)
	}

	return results, nil

}

//query mssql price general
func (r *DataRepository) QueryGeneralPriceMssql() ([]model.DataGeneralPrice, error) {

	rows, err := r.store.dbMssql.Query(r.store.config.Spec.Queryies.GeneralPrice)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	defer rows.Close()

	results := []model.DataGeneralPrice{} // creating empty slice

	for rows.Next() {

		data := &model.DataGeneralPrice{} // creating new struct for every row

		err := rows.Scan(
			&data.Товар,
			&data.ВариантСборки,
			&data.ВариантСборкиРазвернутый,
			&data.Цена,
			&data.СтавкаНДС,
			&data.НДС,
			&data.НачалоДействия,
		)
		if err != nil {
			logger.ErrorLogger.Println(err)
			return nil, err
		}
		results = append(results, *data)
	}

	return results, nil

}

//query mssql sprav
func (r *DataRepository) QuerySprav() ([]model.DataSprav, error) {

	rows, err := r.store.dbMssql.Query(r.store.config.Spec.Queryies.Sprav)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	defer rows.Close()

	results := []model.DataSprav{} // creating empty slice

	for rows.Next() {

		data := &model.DataSprav{} // creating new struct for every row

		err := rows.Scan(
			&data.Наименование,
			&data.НомерСогласноКД,
			&data.Дивизион,
			&data.СтатусМоделиВПроизводстве,
			&data.МассаСнагрузкой,
			&data.МассаБезНагрузки,
			&data.ОписаниеДляПрайса,
			&data.База,
			&data.БазаАвтомобиляДлина,
			&data.ТипКузова,
			&data.ТипФургона,
			&data.ОбозначениеДвигателя,
			&data.ОбъемДвигателя,
			&data.ВидТоплива,
			&data.СтабилизаторЗаднейПодвески,
			&data.ГорныйТормоз,
			&data.ТормознаяСистемаТип,
			&data.ЦветаДопустимыеВЭтомМесяце,
			&data.ОпцииДопустимыеВЭтомМесяце,
			&data.ОпцииПоУмолчанию,
			&data.ЧислоПосадочныхМест,
			&data.ЭкКласс,
			&data.Привод,
			&data.Семейство,
			&data.Лебедка,
			&data.КПП,
			&data.ГБО,
			&data.Надстройка,
			&data.ОсобенностьНадстройки,
			&data.БазовыйТовар,
			&data.ОпцииАЗ,
			&data.ХарактеристикиНоменклатуры,
		)
		if err != nil {
			logger.ErrorLogger.Println(err)
			return nil, err
		}
		results = append(results, *data)
	}

	return results, nil

}

//query mssql options data
func (r *DataRepository) QueryOptionsData() ([]model.DataOptions, error) {

	rows, err := r.store.dbMssql.Query(r.store.config.Spec.Queryies.Options)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	defer rows.Close()

	results := []model.DataOptions{} // creating empty slice

	for rows.Next() {

		data := &model.DataOptions{} // creating new struct for every row

		err := rows.Scan(
			&data.НоменклатураИд,
			&data.НоменклатураНаименование,
			&data.ГруппаОпций,
			&data.ГруппаОпцийНаименование,
			&data.ОпцияИд,
			&data.КраткоеНаименованиеОпции,
			&data.НаименованиеОпции,
			&data.ЗначениеОпцииИд,
			&data.Цена,
			&data.КраткоеНаименование,
			&data.НаименованиеЗначенияОпции,
			&data.Обязательная,
			&data.ВыбранаПоУмолчанию,
			&data.ЭтоПакет,
		)
		if err != nil {
			logger.ErrorLogger.Println(err)
			return nil, err
		}
		results = append(results, *data)
	}

	return results, nil

}

//query mssql options sprav data
func (r *DataRepository) QueryOptionsDataSprav() ([]model.DataOptionsSprav, error) {

	rows, err := r.store.dbMssql.Query(r.store.config.Spec.Queryies.OptionsSprav)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	defer rows.Close()

	results := []model.DataOptionsSprav{} // creating empty slice

	for rows.Next() {

		data := &model.DataOptionsSprav{} // creating new struct for every row

		err := rows.Scan(
			&data.НоменклатураИд,
			&data.НоменклатураНаименование,
			&data.ЗначениеОпции1,
			&data.ЗначениеОпции2,
			&data.КодОпции1,
			&data.КодОпции2,
			&data.ВидСочетания,
		)
		if err != nil {
			logger.ErrorLogger.Println(err)
			return nil, err
		}
		results = append(results, *data)
	}

	return results, nil

}

//query mssql packets data
func (r *DataRepository) QueryPacketsData() ([]model.DataPackets, error) {

	rows, err := r.store.dbMssql.Query(r.store.config.Spec.Queryies.Packets)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	defer rows.Close()

	results := []model.DataPackets{} // creating empty slice

	for rows.Next() {

		data := &model.DataPackets{} // creating new struct for every row

		err := rows.Scan(
			&data.НоменклатураИд,
			&data.НоменклатураНаименование,
			&data.Пакет,
			&data.НаименованиеПакета,
			&data.Опция,
			&data.ЗначениеОпции,
			&data.ЗначениеОпцииНаим,
			&data.ЗначениеОпцииКраткоеНаим,
			&data.ОпцияНаим,
			&data.ОпцияКраткоеНаим,
		)
		if err != nil {
			logger.ErrorLogger.Println(err)
			return nil, err
		}
		results = append(results, *data)
	}

	return results, nil

}

//query mssql colors data
func (r *DataRepository) QueryColorsData() ([]model.DataColors, error) {

	rows, err := r.store.dbMssql.Query(r.store.config.Spec.Queryies.Colors)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	defer rows.Close()

	results := []model.DataColors{} // creating empty slice

	for rows.Next() {

		data := &model.DataColors{} // creating new struct for every row

		err := rows.Scan(
			&data.НоменклатураИд,
			&data.НоменклатураНаименование,
			&data.ЦветИд,
			&data.Наименование,
			&data.ПолноеНаименование,
			&data.ЦветRGB,
			&data.Слойность,
		)
		if err != nil {
			logger.ErrorLogger.Println(err)
			return nil, err
		}
		results = append(results, *data)
	}

	return results, nil

}

//call microservice mailing
func (r *DataRepository) CallMSMailing(data model.DataBooking, config *model.Service) (string, error) {

	bodyBytesReq, err := json.Marshal(data)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return "", err
	}

	resp, err := http.Post(config.Spec.Client.UrlMailingService, "application/json", bytes.NewBuffer(bodyBytesReq))
	if err != nil {
		logger.ErrorLogger.Println(err)
		return "", err
	}

	defer resp.Body.Close()

	bodyBytesResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytesResp), nil

}

//request GAZ CRM booking
func (r *DataRepository) RequestGazCrmApiBooking(data model.DataBooking, config *model.Service) (*model.ResponseGazCrm, error) {

	var dataset model.DataGazCrm
	var response *model.ResponseGazCrm

	bodyJson0 := &model.DataGazCrmReq{
		TimeRequest: data.TimeRequest,
	}
	bodyJson1 := &model.DataGazCrmReq{
		RequestId: data.RequestId,
	}
	bodyJson2 := &model.DataGazCrmReq{
		SubdivisionsId: data.SubdivisionsId,
	}
	bodyJson3 := &model.DataGazCrmReq{
		SubdivisionsName: data.SubdivisionsName,
	}
	bodyJson4 := &model.DataGazCrmReq{
		FormName: data.FormName,
	}
	bodyJson5 := &model.DataGazCrmReq{
		FormId: data.FormId,
	}
	bodyJson6 := &model.DataGazCrmReq{
		HostName: data.HostName,
	}
	bodyJson7 := &model.DataGazCrmReq{
		Division: data.Division,
	}
	bodyJson8 := &model.DataGazCrmReq{
		Area: data.Area,
	}
	bodyJson9 := &model.DataGazCrmReq{
		BrandName: data.BrandName,
	}
	bodyJson10 := &model.DataGazCrmReq{
		CarModel: data.CarModel,
	}
	bodyJson11 := &model.DataGazCrmReq{
		ClientID: data.Clientid,
	}
	bodyJson12 := &model.DataGazCrmReq{
		MetricsType: data.MetricsType,
	}
	bodyJson13 := &model.DataGazCrmReq{
		СlientIP: data.СlientIP,
	}
	bodyJson14 := &model.DataGazCrmReq{
		TypeClient: data.TypeClient,
	}
	bodyJson15 := &model.DataGazCrmReq{
		CompanyName: data.CompanyName,
	}
	bodyJson16 := &model.DataGazCrmReq{
		СlientName: data.Name,
	}
	bodyJson17 := &model.DataGazCrmReq{
		ClientEmail: data.Email,
	}
	bodyJson18 := &model.DataGazCrmReq{
		ClientPhoneNumber: data.PhoneNumber,
	}
	bodyJson19 := &model.DataGazCrmReq{
		Commentary: data.Comment,
	}
	bodyJson20 := &model.DataGazCrmReq{
		AgreementMailing: data.Consentmailing,
	}

	dataset.Data = append(dataset.Data, bodyJson0, bodyJson1, bodyJson2, bodyJson3, bodyJson4, bodyJson5, bodyJson6,
		bodyJson7, bodyJson8, bodyJson9, bodyJson10, bodyJson11, bodyJson12, bodyJson13, bodyJson14,
		bodyJson15, bodyJson16, bodyJson17, bodyJson18, bodyJson19, bodyJson20)

	//d_spaces, err := json.MarshalIndent(dataset, "", "    ")
	//if err != nil {
	//return "", err
	//}

	bodyBytesReq, err := json.Marshal(dataset)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	resp, err := http.Post(config.Spec.Client.UrlGazCrmTest, "application/json", bytes.NewBuffer(bodyBytesReq))
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	defer resp.Body.Close()

	bodyBytesResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	if err := json.Unmarshal(bodyBytesResp, &response); err != nil {
		return nil, err
	}

	return response, nil

}

//request GAZ CRM forms
func (r *DataRepository) RequestGazCrmApiForms(data model.DataForms, config *model.Service) (*model.ResponseGazCrm, error) {

	var dataset model.DataGazCrm
	var response *model.ResponseGazCrm

	bodyJson0 := &model.DataGazCrmReq{
		TimeRequest: data.TimeRequest,
	}
	bodyJson1 := &model.DataGazCrmReq{
		RequestId: data.RequestId,
	}
	bodyJson2 := &model.DataGazCrmReq{
		SubdivisionsId: data.SubdivisionsId,
	}
	bodyJson3 := &model.DataGazCrmReq{
		SubdivisionsName: data.SubdivisionsName,
	}
	bodyJson4 := &model.DataGazCrmReq{
		FormName: data.FormName,
	}
	bodyJson5 := &model.DataGazCrmReq{
		FormId: data.FormId,
	}
	bodyJson6 := &model.DataGazCrmReq{
		HostName: data.HostName,
	}
	bodyJson7 := &model.DataGazCrmReq{
		Division: data.Division,
	}
	bodyJson8 := &model.DataGazCrmReq{
		Area: data.Area,
	}
	bodyJson9 := &model.DataGazCrmReq{
		BrandName: data.BrandName,
	}
	bodyJson10 := &model.DataGazCrmReq{
		CarModel: data.CarModel,
	}
	bodyJson11 := &model.DataGazCrmReq{
		ClientID: data.Clientid,
	}
	bodyJson12 := &model.DataGazCrmReq{
		MetricsType: data.MetricsType,
	}
	bodyJson13 := &model.DataGazCrmReq{
		СlientIP: data.СlientIP,
	}
	bodyJson14 := &model.DataGazCrmReq{
		TypeClient: data.TypeClient,
	}
	bodyJson15 := &model.DataGazCrmReq{
		CompanyName: data.CompanyName,
	}
	bodyJson16 := &model.DataGazCrmReq{
		СlientName: data.Name,
	}
	bodyJson17 := &model.DataGazCrmReq{
		ClientEmail: data.Email,
	}
	bodyJson18 := &model.DataGazCrmReq{
		ClientPhoneNumber: data.PhoneNumber,
	}
	bodyJson19 := &model.DataGazCrmReq{
		Commentary: data.Comment,
	}
	bodyJson20 := &model.DataGazCrmReq{
		AgreementMailing: data.Consentmailing,
	}

	dataset.Data = append(dataset.Data, bodyJson0, bodyJson1, bodyJson2, bodyJson3, bodyJson4, bodyJson5, bodyJson6,
		bodyJson7, bodyJson8, bodyJson9, bodyJson10, bodyJson11, bodyJson12, bodyJson13, bodyJson14,
		bodyJson15, bodyJson16, bodyJson17, bodyJson18, bodyJson19, bodyJson20)

	bodyBytesReq, err := json.Marshal(dataset)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(config.Spec.Client.UrlGazCrmTest, "application/json", bytes.NewBuffer(bodyBytesReq))
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	defer resp.Body.Close()

	bodyBytesResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	if err := json.Unmarshal(bodyBytesResp, &response); err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	return response, nil

}

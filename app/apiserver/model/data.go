package model

import (
	"database/sql/driver"

	validation "github.com/go-ozzo/ozzo-validation"
)

//Data booking
type DataBooking struct {
	RequestId             string `json:"request_id"`
	ActionType            string `json:"action_type"`
	UniqModCode           int    `json:"uniq_mod_code"`
	Modification          string `json:"modification"`
	ModFamily             string `json:"mod_family"`
	ModBodyType           string `json:"mod_body_type"`
	ModEngine             string `json:"mod_engine"`
	ModBase               string `json:"mod_base"`
	ModTuning             string `json:"mod_tuning"`
	Vin                   string `json:"vin"`
	PriceWithNds          int    `json:"price"`
	TypeClient            string `json:"client_type"`
	Inn                   string `json:"inn"`
	Kpp                   string `json:"kpp"`
	Ogrn                  string `json:"ogrn"`
	YurAddressCode        string `json:"reg_address_code"`
	DeliveryAddressCode   string `json:"delivery_address_code"`
	DeliveryAddress       string `json:"delivery_address"`
	Hid                   string `json:"hid"`
	CompanyName           string `json:"client_company_name"`
	RepresentativeName    string `json:"representative_name"`
	RepresentativeSurname string `json:"representative_surname"`
	Surname               string `json:"surname"`
	Name                  string `json:"client_name"`
	Patronymic            string `json:"patronymic"`
	PassportSer           string `json:"passport_ser"`
	PassportNumber        string `json:"passport_number"`
	Snils                 string `json:"snils"`
	DateOfBirth           string `json:"date_of_birth"`
	Email                 string `json:"client_email"`
	PhoneNumber           string `json:"client_phone_number"`
	Comment               string `json:"commentary"`
	Consentmailing        string `json:"agreement_mailing"`
	TimeRequest           string `json:"event_datetime"`
	File                  string `json:"file"`
	BillNumber            string `json:"bill_namber"`
	UrlMod                string `json:"url_mod"`
	Clientid              string `json:"clientid_google"` //Google Analytics cookies
	Ymuid                 string `json:"ClientID"`        //Yandex Metrics cookies
	TestMod               bool   `json:"testmod"`         //true - test, false - prod
	//fields for gaz crm
	SubdivisionsId   string `json:"subdivisions_id"`
	SubdivisionsName string `json:"subdivisions_name"`
	FormName         string `json:"form_name"`
	FormId           string `json:"id_form"`
	HostName         string `json:"host_name"`
	Division         string `json:"division"`
	Area             string `json:"area"`
	BrandName        string `json:"brand_name"`
	CarModel         string `json:"car_model"`
	MetricsType      string `json:"metrics_type"`
	СlientIP         string `json:"client_IP"`
}

//Validation data booking
func (d *DataBooking) ValidateDataBooking() error {
	return validation.ValidateStruct(
		d,
		validation.Field(&d.RequestId, validation.Required),
		validation.Field(&d.UniqModCode, validation.Required),
		validation.Field(&d.Modification, validation.Required),
		validation.Field(&d.ModFamily, validation.Required),
		validation.Field(&d.ModBodyType, validation.Required),
		validation.Field(&d.ModEngine, validation.Required),
		validation.Field(&d.ModBase, validation.Required),
		validation.Field(&d.ModTuning, validation.Required),
		validation.Field(&d.Vin, validation.Required),
		validation.Field(&d.PriceWithNds, validation.Required),
		validation.Field(&d.TypeClient, validation.In("company", "personal")),
		validation.Field(&d.Inn, validation.Required),
		validation.Field(&d.Kpp, validation.Required),
		validation.Field(&d.Ogrn, validation.Required),
		validation.Field(&d.YurAddressCode, validation.Required),
		validation.Field(&d.DeliveryAddressCode, validation.Required),
		validation.Field(&d.Hid, validation.Required),
		validation.Field(&d.CompanyName, validation.Required),
		validation.Field(&d.Surname, validation.Required),
		validation.Field(&d.Name, validation.Required),
		validation.Field(&d.Patronymic, validation.Required),
		validation.Field(&d.DateOfBirth, validation.Date("2006-01-02")),
		validation.Field(&d.PassportSer, validation.Required),
		validation.Field(&d.PassportNumber, validation.Required),
		validation.Field(&d.Email, validation.Required),
		validation.Field(&d.PhoneNumber, validation.Required),
		validation.Field(&d.TimeRequest, validation.Date("2006-01-02T15:04:05")),
		validation.Field(&d.Patronymic, validation.Required),
		validation.Field(&d.Consentmailing, validation.In("yes", "no")),
		//validation fo gaz crm fields
		validation.Field(&d.Division, validation.In("lcv/mcv", "bus")),
		validation.Field(&d.Area, validation.In("dealer", "distrib")),
		validation.Field(&d.MetricsType, validation.In("yandex")),
		validation.Field(&d.TypeClient, validation.In("company", "personal")),
		validation.Field(&d.TimeRequest, validation.Date("2006-01-02T15:04:05")),
		validation.Field(&d.ActionType, validation.Required, validation.In("form", "bill", "acquiring")),
	)
}

//sql null handling
type NullString string

func (s *NullString) Scan(value interface{}) error {
	if value == nil {
		*s = "nil"
		return nil
	}
	strVal, ok := value.(string)
	if !ok {
		return nil
	}
	*s = NullString(strVal)
	return nil
}

func (s NullString) Value() (driver.Value, error) {
	if len(s) == 0 { // if nil or empty string
		return nil, nil
	}
	return string(s), nil
}

//Data get
type DataStocks struct {
	VIN                            string
	Площадка                       string
	Наименование_номенклатуры      string
	Номер_согласно_КД              string
	Дивизион                       string
	Доработчик_Подрядчик           NullString //NULL
	Test_truck                     bool
	Телематика                     string
	Номер_шасси                    string
	Номер_двигателя                NullString //NULL
	Грузоподъемность_кг            string
	Цвет                           string
	Вариант_сборки                 string
	Расшифровка_варианта_сборки    string
	Вариант_сборки_свернутый       NullString //NULL
	Год_VIN                        string
	Дата_сборки                    NullString //NULL
	Справочная_стоимость_по_прайсу string
	Дата_отгрузки                  NullString //NULL
	Дата_прихода                   NullString //NULL
	Страна                         NullString //NULL
	Контрагент_получателя          string
	Стоянка                        string
	Город_стоянки                  NullString //NULL
	Площадка_получателя_Ид         string
	Контрагент_получателя_Ид       string
	Город_стоянки_Ид               NullString //NULL
	Номер_заявки                   NullString //NULL
	Для_доработки                  NullString //NULL
	Номерной_товар                 string
}

//data price basic models
type DataBasicModelsPrice struct {
	Товар          string
	НачалоДействия string
	Цена           string
	НДС            string
	СтавкаНДС      string
}

//data price options
type DataOptionsPrice struct {
	ЕНСП_Модификация_Ид NullString
	Товар               string
	ЗначениеОпции       string
	ОбозначениеОпции    string
	Цена                string
	СтавкаНДС_Ид        string
	НДС                 string
	НачалоДействия      string
	СоставПакета        NullString
}

//data price general
type DataGeneralPrice struct {
	Товар                    string
	ВариантСборки            string
	ВариантСборкиРазвернутый string
	Цена                     string
	СтавкаНДС                string
	НДС                      string
	НачалоДействия           string
}

//data sprav
type DataSprav struct {
	Наименование               string
	НомерСогласноКД            string
	Дивизион                   string
	СтатусМоделиВПроизводстве  string
	МассаСнагрузкой            string
	МассаБезНагрузки           string
	ОписаниеДляПрайса          string
	База                       string
	БазаАвтомобиляДлина        string
	ТипКузова                  string
	ТипФургона                 string
	ОбозначениеДвигателя       string
	ОбъемДвигателя             string
	ВидТоплива                 string
	СтабилизаторЗаднейПодвески string
	ГорныйТормоз               string
	ТормознаяСистемаТип        string
	ЦветаДопустимыеВЭтомМесяце string
	ОпцииДопустимыеВЭтомМесяце string
	ОпцииПоУмолчанию           string
	ЧислоПосадочныхМест        string
	ЭкКласс                    string
	Привод                     string
	Семейство                  string
	Лебедка                    string
	КПП                        string
	ГБО                        string
	Надстройка                 string
	ОсобенностьНадстройки      string
	БазовыйТовар               NullString
	ОпцииАЗ                    string
	ХарактеристикиНоменклатуры string
}

//options data
type DataOptions struct {
	НоменклатураИд            string
	НоменклатураНаименование  string
	ГруппаОпций               NullString
	ГруппаОпцийНаименование   NullString
	ОпцияИд                   string
	КраткоеНаименованиеОпции  string
	НаименованиеОпции         string
	ЗначениеОпцииИд           string
	Цена                      NullString
	КраткоеНаименование       string
	НаименованиеЗначенияОпции string
	Обязательная              string
	ВыбранаПоУмолчанию        string
	ЭтоПакет                  string
}

//options sprav data
type DataOptionsSprav struct {
	НоменклатураИд           string
	НоменклатураНаименование string
	ЗначениеОпции1           string
	ЗначениеОпции2           string
	КодОпции1                string
	КодОпции2                string
	ВидСочетания             string
}

//packets data
type DataPackets struct {
	НоменклатураИд           string
	НоменклатураНаименование string
	Пакет                    string
	НаименованиеПакета       string
	Опция                    string
	ЗначениеОпции            string
	ЗначениеОпцииНаим        string
	ЗначениеОпцииКраткоеНаим string
	ОпцияНаим                string
	ОпцияКраткоеНаим         string
}

//colors data
type DataColors struct {
	НоменклатураИд           string
	НоменклатураНаименование string
	ЦветИд                   string
	Наименование             string
	ПолноеНаименование       string
	ЦветRGB                  string
	Слойность                string
}

//data options
type OptionsData struct {
}

//Data forms
type DataForms struct {
	//gaz crm fields
	TimeRequest      string `json:"event_datetime"` //general field with booking
	RequestId        string `json:"request_id"`     //general field with booking
	SubdivisionsId   string `json:"subdivisions_id"`
	SubdivisionsName string `json:"subdivisions_name"`
	FormName         string `json:"form_name"`
	FormId           string `json:"id_form"`
	HostName         string `json:"host_name"`
	Division         string `json:"division"`
	Area             string `json:"area"`
	BrandName        string `json:"brand_name"`
	CarModel         string `json:"car_model"`
	Clientid         string `json:"ClientID"` //general field with booking
	MetricsType      string `json:"metrics_type"`
	СlientIP         string `json:"client_IP"`
	TypeClient       string `json:"client_type"`         //general field with booking
	CompanyName      string `json:"client_company_name"` //general field with booking
	Name             string `json:"client_name"`         //general field with booking
	Email            string `json:"client_email"`        //general field with booking
	PhoneNumber      string `json:"client_phone_number"` //general field with booking
	Comment          string `json:"commentary"`          //general field with booking
	Consentmailing   string `json:"agreement_mailing"`   //general field with booking
	//additional fields
	ActionType   string `json:"action_type"`
	Modification string `json:"modification"`
	ModFamily    string `json:"mod_family"`
	ModBodyType  string `json:"mod_body_type"`
	ModEngine    string `json:"mod_engine"`
	ModBase      string `json:"mod_base"`
	ModTuning    string `json:"mod_tuning"`
	Vin          string `json:"vin"`
	PriceWithNds int    `json:"price"`
	UrlMod       string `json:"url_mod"`
}

//Validation data fiz
func (d *DataForms) ValidateDataForms() error {
	return validation.ValidateStruct(
		d,
		validation.Field(&d.RequestId, validation.Required),
		validation.Field(&d.SubdivisionsId, validation.Required),
		validation.Field(&d.SubdivisionsName, validation.Required),
		validation.Field(&d.FormName, validation.Required),
		validation.Field(&d.FormId, validation.Required),
		validation.Field(&d.HostName, validation.Required),
		validation.Field(&d.BrandName, validation.Required),
		validation.Field(&d.CarModel, validation.Required),
		validation.Field(&d.Clientid, validation.Required),
		validation.Field(&d.MetricsType, validation.Required),
		validation.Field(&d.Name, validation.Required),
		validation.Field(&d.Email, validation.Required),
		validation.Field(&d.PhoneNumber, validation.Required),
		validation.Field(&d.UrlMod, validation.Required),
		validation.Field(&d.Modification, validation.Required),
		validation.Field(&d.ModFamily, validation.Required),
		validation.Field(&d.ModBodyType, validation.Required),
		validation.Field(&d.ModEngine, validation.Required),
		validation.Field(&d.ModBase, validation.Required),
		validation.Field(&d.ModTuning, validation.Required),
		validation.Field(&d.Vin, validation.Required),
		validation.Field(&d.PriceWithNds, validation.Required),
		validation.Field(&d.Division, validation.In("lcv/mcv", "bus")),
		validation.Field(&d.Area, validation.In("dealer", "distrib")),
		validation.Field(&d.MetricsType, validation.In("yandex")),
		validation.Field(&d.TypeClient, validation.In("company", "personal")),
		validation.Field(&d.TimeRequest, validation.Date("2006-01-02T15:04:05")),
		validation.Field(&d.ActionType, validation.Required, validation.In("form")),
		validation.Field(&d.Consentmailing, validation.In("yes", "no")),
	)
}

//gaz crm
//data struct for call gaz crm api method
type DataGazCrm struct {
	Data DataGazCrmReq `json:"Data"`
}

//data struct for call gaz crm api method
type DataGazCrmReq struct {
	//gaz crm fields
	TimeRequest struct {
		TimeRequest string `json:"event_datetime"` //general field with booking
	}
	RequestId struct {
		RequestId string `json:"request_id"` //general field with booking
	}
	SubdivisionsId struct {
		SubdivisionsId string `json:"subdivisions_id"`
	}
	SubdivisionsName struct {
		SubdivisionsName string `json:"subdivisions_name"`
	}
	FormName struct {
		FormName string `json:"form_name"`
	}
	FormId struct {
		FormId string `json:"id_form"`
	}
	HostName struct {
		HostName string `json:"host_name"`
	}
	Division struct {
		Division string `json:"division"`
	}
	Area struct {
		Area string `json:"area"`
	}
	BrandName struct {
		BrandName string `json:"brand_name"`
	}
	CarModel struct {
		CarModel string `json:"car_model"`
	}
	ClientID struct {
		ClientID string `json:"ClientID"`
	}
	MetricsType struct {
		MetricsType string `json:"metrics_type"`
	}
	СlientIP struct {
		СlientIP string `json:"client_IP"`
	}
	TypeClient struct {
		TypeClient string `json:"client_type"` //general field with booking
	}
	CompanyName struct {
		CompanyName string `json:"client_company_name"` //general field with booking
	}
	СlientName struct {
		СlientName string `json:"client_name"` //general field with booking
	}
	ClientEmail struct {
		ClientEmail string `json:"client_email"` //general field with booking
	}
	ClientPhoneNumber struct {
		ClientPhoneNumber string `json:"client_phone_number"` //general field with booking
	}
	Commentary struct {
		Commentary string `json:"commentary"` //general field with booking
	}
	AgreementMailing struct {
		AgreementMailing string `json:"agreement_mailing"` //general field with booking
	}
}

//data struct for call gaz crm api method
//lead_get gaz crm
type DataLeadGet struct {
	Data DataLeadGet_Gazcrm `json:"Data"`
}

//lead_get gaz crm
type DataLeadGet_Gazcrm struct {
	TimeRequest struct {
		TimeRequest string `json:"event_datetime"`
	}
	EventName struct {
		EventName string `json:"event_name"`
	}
	RequestId struct {
		RequestId string `json:"request_id,omitempty"`
	}
	SubdivisionsId struct {
		SubdivisionsId string `json:"subdivisions_id,omitempty"`
	}
	SubdivisionsName struct {
		SubdivisionsName string `json:"subdivisions_name"`
	}
	FormName struct {
		FormName string `json:"form_name"`
	}
	HostName struct {
		HostName string `json:"host_name"`
	}
	Division struct {
		Division string `json:"division"`
	}
	Area struct {
		Area string `json:"area"`
	}
	BrandName struct {
		BrandName string `json:"brand_name"`
	}
	ClientID struct {
		ClientID string `json:"ClientID"`
	}
	MetricsType struct {
		MetricsType string `json:"metrics_type"`
	}
}

//work_list gaz crm
type DataWorkList struct {
	Data DataWorkList_Gazcrm `json:"Data"`
}

//work_list gaz crm
type DataWorkList_Gazcrm struct {
	TimeRequest struct {
		TimeRequest string `json:"event_datetime"`
	}
	EventName struct {
		EventName string `json:"event_name"`
	}
	GazcrmClientId struct {
		GazcrmClientId string `json:"gazcrm_client_id"`
	}
	GazCrmWorkListId struct {
		GazCrmWorkListId string `json:"gazcrm_worklist_id"`
	}
}

//status gaz crm
type DataStatuses struct {
	Data DataStatuses_Gazcrm `json:"Data"`
}

//status gaz crm
type DataStatuses_Gazcrm struct {
	TimeRequest struct {
		TimeRequest string `json:"event_datetime"`
	}
	EventName struct {
		EventName string `json:"event_name"`
	}
	RequestId struct {
		RequestId string `json:"request_id"`
	}
	GazcrmClientId struct {
		GazcrmClientId string `json:"gazcrm_client_id"`
	}
	GazCrmWorkListId struct {
		GazCrmWorkListId string `json:"gazcrm_worklist_id"`
	}
	ClientID struct {
		ClientID string `json:"ClientID"`
	}
	MetricsType struct {
		MetricsType string `json:"metrics_type"`
	}
}

//resp struct api gaz crm
type ResponseGazCrm struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

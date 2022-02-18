package main

/*
type department struct {
	name         string
	responsibles []string
}*/

const triggerHello = "hello"
const triggerDayOff = "izin"

const hintDayOff = "izin almak için lütfen aşağıdaki formatla girişinizi yapınız. \n" +
	"/izin-Al 	İsim	 Soyisim	 İzin başlangıç tarihi		İzin bitiş tarihi \n" +
	"/izin-Al Evren Ispiroglu 06.02.2021 07.02.2021"

const bilisimURL = "https://www.google.com/url?sa=i&url=https%3A%2F%2Fwww.expogi.com%2Fyonetim-bilisim-sistemleri-bilisim-a-s.html%2F&psig=AOvVaw1JKgGLdWec2PbXEeTVojo6&ust=1644442690751000&source=images&cd=vfe&ved=0CAgQjRxqFwoTCICrkciI8fUCFQAAAAAdAAAAABAD"
const bilisimLURL = "https://www.google.com/url?sa=i&url=https%3A%2F%2Ftr.linkedin.com%2Fcompany%2Fbilisim-as&psig=AOvVaw1JKgGLdWec2PbXEeTVojo6&ust=1644442690751000&source=images&cd=vfe&ved=0CAgQjRxqFwoTCICrkciI8fUCFQAAAAAdAAAAABAJ"

/*
   type dayOffProps struct {
   	text       string
   	Fallback   string `json:"fallback"`
   	color      string
   	Pretext    string `json:"pretext"`
   	AuthorName string `json:"author_name"`
   	AuthorIcon string `json:"author_icon"`
   	AuthorLink string `json:"author_link"`
   	Title      string `json:"title"`
   	TitleLink  string `json:"title_link"`
   	Fields     []struct {
   		Short bool   `json:"short"`
   		Title string `json:"title"`
   		Value string `json:"value"`
   	} `json:"fields"`
   	Field struct {
   		Short bool   `json:"short"`
   		Title string `json:"title"`
   		Value string `json:"value"`
   	}
   	ImageURL string `json:"image_url"`
   }

   func dayOffPostInit(off dayOff) map[string]interface{} {
   	props := &dayOffProps{
   		color:      "FF0000",
   		AuthorName: "BilisimHR",
   		AuthorLink: "https://bilisim.com.tr/tr/urunler/bilisim-hr-insan-kaynaklari-yazilimi",
   		Title:      "BilisimHR",
   		TitleLink:  "https://bilisim.com.tr/tr/urunler/bilisim-hr-insan-kaynaklari-yazilimi",
   		text:       off.toString(),
   	}
   	println(props)
   	m := structs.Map(props)
   	return m
   }*/

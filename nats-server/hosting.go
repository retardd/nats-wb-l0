package nats_server

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"l0/internal/datastorage/cache"
	"l0/internal/nats"
	"l0/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type Api struct {
	cch  *cache.Cache
	port int
	http.Server
	pub *nats.Pub
}

func InitApi(cch *cache.Cache, pub *nats.Pub) *Api {
	api := Api{}
	api.cch = cch
	port := utils.GetIntEnv("API_PORT", 8080)
	api.port = port
	api.Addr = fmt.Sprintf(":%d", port)
	api.Handler = api.createRouter()
	api.pub = pub
	return &api
}

// Start запуск сервера
func (api *Api) Start() {
	defer api.Close()
	go func() {
		fmt.Printf("Сервер запущен на http://localhost%s\n", api.Addr)
		log.Fatal(api.ListenAndServe())
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
}

func (api *Api) Close() {
	if err := api.Shutdown(context.Background()); err != nil {
		panic(err)
	}
}

// createRouter создание роутера и настройка обработчиков
func (api *Api) createRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", api.indexHandler)
	mux.HandleFunc("/submit", api.submitHandler)
	mux.HandleFunc("/gofakeit", api.addButtonHandler)

	return mux
}

// Обработчик для главной страницы
func (api *Api) indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "ERROR TEMPLATE PARSE", http.StatusInternalServerError)
	}

	// Отправка HTML на клиент
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "ERROR TEMPLATE EXECUTE", http.StatusInternalServerError)
	}
	return
}

// Обработчик для отправки данных
func (api *Api) submitHandler(w http.ResponseWriter, r *http.Request) {
	// Получение данных из формы
	data := r.FormValue("inputData")

	// Вывод данных в консоль
	fmt.Printf("POST FORM DATA: %s\n", data)

	value, err := strconv.Atoi(data)
	if err != nil {
		log.Fatalf("ERROR ID INT-STRING FORM DATA: %s", err)
	}

	err, mdl := api.cch.FindData(context.TODO(), value)

	var enodata []byte

	if err != nil {
		enodata = []byte("FIND DATA ERROR (CHECK ID)")
		// Отправка ответа обратно на клиент
		_, err = w.Write(enodata)
	} else {
		enodata, _ = json.Marshal(mdl)
		// Отправка ответа обратно на клиент
		_, err = w.Write(enodata)
	}

	if err != nil {
		log.Fatalf("ERROR SUBMIT FORM DATA: %s", err)
	}
	return
}

func (api *Api) addButtonHandler(w http.ResponseWriter, r *http.Request) {
	err, model := utils.FakeModel()

	if err != nil {
		log.Fatalf("ERROR FAKE MODEL CREATOR: %s", err)
	}

	go func() {
		err := api.pub.GetPublish(&model)
		if err != nil {
			log.Fatalf("ERROR ADD DATA TO CACHE: %s", err)
		}
	}()

	type reqData struct {
		Did int
	}

	api.cch.Mtx.RLock()
	byteArr, _ := json.Marshal(reqData{Did: api.cch.LastId})
	api.cch.Mtx.RUnlock()

	// Отправка ответа обратно на клиент
	_, err = w.Write(byteArr)

	if err != nil {
		log.Fatalf("ERROR SUBMIT ID: %s", err)
	}
	return
}

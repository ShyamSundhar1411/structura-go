
# **Structura**  
**Automated Folder Structure & Dependency Management for Go Projects**  

---

## **Overview**  
Structura is a powerful **CLI tool** designed to automate the creation of consistent and production-ready folder structures for **Golang projects**. It supports multiple frameworks, including **Gin, Echo, Fiber, and Chi**, and offers flexible **YAML-based configuration** for managing dependencies and customizing project layouts.  

Whether you're building a RESTful API or a microservice, Structura handles:  
- **Folder structuring**  
- **Dependency installation**  
- **Environment setup**  
- **Custom configurations** via YAML  

---

## **Features**  
- **Multiple Architectures:** MVC, MVCS, Hexagonal, and more  
- **Auto-Generates Folders & Boilerplate Files:** (e.g., `env.go`, `.env`)  
- **Dependency Management:** Installs `viper`, `logrus`, `gin`, `echo`, etc., using `go get`  
- **YAML Configuration:** Flexible and customizable project initialization  
- **Cross-Platform Compatibility:** Works on Windows, Linux, and macOS  

---

## **Installation**  
Install Structura using `go install`:  
```bash
go install github.com/ShyamSundhar1411/structura-go@latest
```  
**Compatibility:** Go `>=1.18`  

---

## **Usage**  

### **Initialize a New Project**  
Create a new Go project with the desired architecture:  
```bash
structura-go init myproject --framework gin
```
For Echo:  
```bash
structura-go init myproject --framework echo
```
For Fiber:  
```bash
structura-go init myproject --framework fiber
```
For Chi:  
```bash
structura-go init myproject --framework chi
```

### **Generate Project Files**  
If you already have a project, you can simply generate the structure:  
```bash
structura-go init
```

---

## **Folder Structure Example**  

When you run `structura-go init`, it generates the following folder structure:  

```plaintext
/myproject
├── bootstrap
│   ├── env.go
├── config
│   ├── config.go
├── internal
│   ├── handlers
│   │   ├── user_handler.go
│   ├── models
│   │   ├── user.go
│   ├── services
│   │   ├── user_service.go
├── .env
├── main.go
├── go.mod
├── go.sum
```  
**Architecture Variations:**  
- `MVC`: `models`, `services`, `handlers`, `routes`  
- `MVCS`: Adds `services` layer for business logic separation  
- `Hexagonal`: Adds `adapters` and `ports` folders for dependency inversion  


## **Contributing**  
We welcome contributions!  
To contribute:  
1. Fork the repository  
2. Create a new feature branch  
3. Commit your changes  
4. Open a Pull Request (PR)  

---

## **License**  
Structura is licensed under the **MIT License**.  
Feel free to use, modify, and distribute it.  

---

## **Feedback & Issues**  
If you encounter any issues or have suggestions, feel free to open an issue on [GitHub](https://github.com/ShyamSundhar1411/structura-go/issues).  

---
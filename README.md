# Structura 🚀  
**Automated Folder Structure & Dependency Management for Go Projects**  

## Overview  
Structura is a project scaffolding tool that **automatically generates folder structures** (MVC, Hexagonal, MVCS) and **installs dependencies** with their required configurations for Go projects using the **Gin** or **Echo** framework.  

## Features  
✅ Supports multiple architectures: **MVC, MVCS, Hexagonal, etc.**  
✅ Auto-generates **folders & boilerplate files** (e.g., `env.go`, `.env`)  
✅ Installs dependencies (`viper`, `logrus`, `gin`, etc.) using `go get`  
✅ Configurable via **YAML files**  
✅ Flexible & extensible  

## Installation  
```bash
go install github.com/ShyamSundhar1411/structura-go@latest
```

## Usage  
### Initialize a new project  
```bash
structura init myproject --framework gin
```


### Generate project files  
```bash
structura init
```

## Folder Structure Example  
```
myproject/
│── bootstrap/
│   ├── env.go
│── .env
│── main.go
│── go.mod
```

## Contributing  
1. Fork the repo  
2. Create a feature branch  
3. Open a PR 🚀  

## License  
MIT License  

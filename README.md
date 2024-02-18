# API de Pagos

La API de Pagos es un sistema diseñado para manejar transacciones online, permitiendo a los comerciantes procesar pagos, consultar detalles de pagos anteriores y gestionar reembolsos.
![alt text](https://github.com/davidgg090/paymentApiGO/blob/main/img/payment.png?raw=true)

## Comenzando

Estas instrucciones te proporcionarán una copia del proyecto en funcionamiento en tu máquina local para fines de desarrollo y pruebas.

### Prerrequisitos

Qué cosas necesitas para instalar el software y cómo instalarlas:

- Go (versión 1.15 o superior)
- PostgreSQL
- [Opcional] Docker, si prefieres ejecutar la base de datos en un contenedor

### Instalación

Un paso a paso que te guiará sobre cómo poner en marcha un entorno de desarrollo:

1. Clona el repositorio en tu máquina local:

```bash
git clone https://github.com/davidgg090/tuProyecto.git
```

2. Accede al directorio del proyecto:

```bash
cd tuProyecto
```

3. Copia el archivo de configuración de ejemplo:

```bash
cp .env.example .env
```

4. Edita el archivo `.env` con tus credenciales de base de datos

5. Instala las dependencias del proyecto:

```bas
  go mod tidy
```

6. Inicia la base de datos PostgreSQL y crea la base de datos requerida.
7. Ejecuta el proyecto
    
```bash
go run cmd/payment/main.go
```


### Uso

Para empezar a usar la API, puedes importar la colección de Postman proporcionada en la carpeta `postman` del repositorio. Esta colección incluye ejemplos de solicitudes para todos los endpoints disponibles.


### Uso con Docker

1. Construir la imagen Docker:
    
    ```bash
    docker build -t api-pagos .
    ```

2. Ejecutar el contenedor:

    ```bash
    docker run -dp 8080:8080 --env-file .env -v $(pwd)/.env:/.env payment-api
    ```
   


## Licencia

Este proyecto está bajo la Licencia MIT


    
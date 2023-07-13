# Photo Bomb

## Project Description
Through the use of image steganography, our website allows users to hide text in their photos. Users are able to create their own accounts where they upload images with embeded messages. These message can allow for keeping a private journal, hiding information in plain sight, and creating hidden watermarks.

## Group Members

### Frontend
Aszti Chadzynski

Eric Truong
### Backend
Tanner Garcia

Dominik Kapuscinski


## Requirements
- Go (project built on 1.19.5)
- Angular
- MySQL

## Usage
First clone the project locally
```
git clone https://github.com/Dominikk7/PhotoBomb
```
Next, start a MySQL instance.

Ensure that 2 databases exist: images and users

Create directory "uploads" in the backend directory

Navigate to the backend/cmd directory and start the backend go server
```
go run main.go
```
In a separate terminal, navigate to the frontend folder. 

Build and run the angular application
```
npm install
ng serve
```
The application is now running.  Navigate to the provided link in your web browser to view it.

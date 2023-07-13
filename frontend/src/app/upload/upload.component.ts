import { Component, Injectable } from '@angular/core';
import { HttpClient,HttpHeaders } from '@angular/common/http';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { DecodeComponent } from '../decode/decode.component';
import { EncodeComponent } from '../encode/encode.component';
import { Observable } from 'rxjs/internal/Observable';
import { imageObj } from '../album/imageObj';
import { ipAddress } from '../../../ip.conf';

@Injectable({
  providedIn: 'root'
})

@Component({
  selector: 'app-upload',
  templateUrl: './upload.component.html',
  styleUrls: ['./upload.component.css']
})


export class UploadComponent {
  

  fileName = '';
  ipAddy = ipAddress.ipFull;
  formData = new FormData();
  hide = true;
  message!: string;
  public displayImageList= [""];
  public imageTest!: imageObj[]

  

  constructor(private httpClient: HttpClient, private dialog: MatDialog) {
    hide: true;
    console.log("new instance");
  }

  openEncode() {
    const dialogConfig = new MatDialogConfig();
    dialogConfig.disableClose = true;

    let dialogRef = this.dialog.open(EncodeComponent, dialogConfig);
    dialogRef.componentInstance.service = this;
  }

  openDecode() {
    const dialogConfig = new MatDialogConfig();
    dialogConfig.disableClose = true;

    let dialogRef = this.dialog.open(DecodeComponent, dialogConfig);
    dialogRef.componentInstance.service = this;
  }

  getDecoded(){
    console.log("Sent");
    return this.httpClient.post('http://' + this.ipAddy +'/upload/decode', this.formData, //8080
    {
      withCredentials: true,
      
    });
  }

  getEncoded(){
    console.log("Sent");
    return this.httpClient.post('http://' + this.ipAddy +'/upload/encode', this.formData, //8080
    {
      withCredentials: true
    });
  }

  

  onFileSelected(event: Event) {
    const target= event.target as HTMLInputElement;
    const file: File = (target.files as FileList)[0];
    
    if (file) {
      this.fileName = file.name;    
      this.formData = new FormData(); 
      this.formData.append("uploadfile", file);
      //this.formData.append("imagetext", "imageText");         
    }
    
  }
  
  
}



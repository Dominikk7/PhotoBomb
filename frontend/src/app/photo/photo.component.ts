import { Component, Inject, OnInit, TemplateRef } from '@angular/core';
import { AlbumComponent } from '../album/album.component';
import { MAT_DIALOG_DATA, MatDialog, MatDialogConfig, MatDialogRef } from '@angular/material/dialog';
import { HttpClient,HttpHeaders } from '@angular/common/http';
import { AlbumDecodeComponent } from '../album-decode/album-decode.component';
import { ExistDecodeComponent } from '../exist-decode/exist-decode.component';
import { ComponentType } from '@angular/cdk/portal';
import { ipAddress } from '../../../ip.conf';

@Component({
  selector: 'app-photo',
  templateUrl: './photo.component.html',
  styleUrls: ['./photo.component.css']
})
export class PhotoComponent implements OnInit{
  public service!: AlbumComponent;
  photoURL! : string;
  timestamp!: string;
  ipAddy = ipAddress.ipFull;
  imageText= "";
  hide = true;


  constructor(
    public dialog: MatDialog,
    private httpClient: HttpClient,
    public dialogPhoto: MatDialogRef<PhotoComponent>, public dialogDec: MatDialog
    ){}
    
  ngOnInit(): void {
    this.Init();
  }

  Init(){
    this.photoURL = this.service.photoURL;
    this.timestamp = this.photoURL.substring(this.photoURL.lastIndexOf("=")+1);
  }

    
  Decode1(){
    this.httpClient.get(this.ipAddy+"/decode/?timestamp="+this.timestamp,
      {
        withCredentials: true
      }).subscribe((rslt:Object)=>{
        this.imageText = JSON.stringify(rslt);
        console.log(this.imageText);
      }) ;

    this.hide = !this.hide;

  }


    Decode(){
      const dialogConfig = new MatDialogConfig();
      dialogConfig.disableClose = true;
  
      let dialogRef = this.dialogDec.open(AlbumDecodeComponent, dialogConfig);
      dialogRef.componentInstance.timestamp = this.timestamp;
      dialogRef.componentInstance.Init();
    }

    Delete(){
      console.log("delete");
      this.httpClient.delete(this.ipAddy+"/delete/?timestamp="+this.timestamp,
      {
        withCredentials: true
      }).subscribe((rslt:Object)=>{
        console.log(JSON.stringify(rslt))
      }) ;

      this.dialogPhoto.close();
    }

    Update(){
      this.service.removeImage(this.photoURL);
    }

    Download(){
      console.log("download"); 
      window.open(this.photoURL);
  }

    Close(){
      this.dialogPhoto.close();
    }

  

}

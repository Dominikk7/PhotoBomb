import { Component } from '@angular/core';
import { DomSanitizer, SafeUrl } from '@angular/platform-browser';
import { HttpClient,HttpHeaders } from '@angular/common/http';
import { imageObj } from './imageObj';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { PhotoComponent } from '../photo/photo.component';
import { ipAddress } from '../../../ip.conf';

@Component({
  selector: 'app-album',
  templateUrl: './album.component.html',
  styleUrls: ['./album.component.css']
})
export class AlbumComponent {
  fileName = '';
  ipAddy = ipAddress.ipFull;//localhost
  formData = new FormData();
  hide = true;
  message!: string;
  public displayImageList= [""];
  public imageTest!: imageObj[]
  photoURL! : string
  

  constructor(private httpClient: HttpClient, private readonly domSanitizer: DomSanitizer, private dialog: MatDialog) {
    hide: true;
    console.log("new instance");
  }


  openPhoto() {
    const dialogConfig = new MatDialogConfig();

    let dialogRef = this.dialog.open(PhotoComponent, {disableClose: true, data: this.photoURL, 
    width: '600px',});
    dialogRef.componentInstance.service = this;
  }

  ngOnInit() {
    this.updateImages();
  }

  getImageList(){
    console.log("Sent");
    return this.httpClient.get('http://' + this.ipAddy +'/download/list/',
    {
      withCredentials: true
    });
  }

  updateImages(){
    this.displayImageList.splice(0, this.displayImageList.length);
    this.getImageList().subscribe(images => {
      this.imageTest = JSON.parse(JSON.stringify(images));;
      for(let image in this.imageTest){
        var url = "http://"+ this.ipAddy + "/download/?timestamp="
        var timestamp = this.imageTest[image].timestamp;
        this.displayImageList.push(url + timestamp)
      }
     });
    this.displayImageList.pop();

    console.log(this.displayImageList);
  }

  removeImage(url:string){
    const index = this.displayImageList.indexOf(url);
    if (index !== -1) {
      this.displayImageList.splice(index, 1);
  }     
  }

  public highlight(ev: MouseEvent) {  
    const target = ev.target as HTMLImageElement;
    console.log(target.currentSrc);
    this.photoURL = target.currentSrc;
    if (this.photoURL.length > 0){
      this.openPhoto();
    }
    
  }

}

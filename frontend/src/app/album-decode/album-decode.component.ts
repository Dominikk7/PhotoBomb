import { Component, OnInit } from '@angular/core';
import { PhotoComponent } from '../photo/photo.component';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { HttpClient,HttpHeaders } from '@angular/common/http';
import { ipAddress } from '../../../ip.conf';

@Component({
  selector: 'app-album-decode',
  templateUrl: './album-decode.component.html',
  styleUrls: ['./album-decode.component.css']
})
export class AlbumDecodeComponent implements OnInit{
  public decodeComp!: PhotoComponent;
  currentCustomer = 'Maria';
  public test="Test";
  timestamp!: string;
  loaded!: boolean;

  constructor(
    public dialog: MatDialogRef<AlbumDecodeComponent>, private httpClient: HttpClient,
    ){} 
  ngOnInit(): void {
    //console.log(this.timestamp + "test");
  }
    
    Init() {
      console.log(this.timestamp + "test");

      this.httpClient.get(ipAddress.ipFull +"/decode/?timestamp="+this.timestamp,
      {
        withCredentials: true
      }).subscribe((rslt:Object)=>{
        this.test = JSON.stringify(rslt);
        console.log(this.test);
      }) ;
    }

  close() {
    console.log(this.timestamp);
    this.dialog.close();
  }
}

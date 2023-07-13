import { HttpClient } from '@angular/common/http';
import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';

import { StorageService } from '../authweb/storage.service';
import { CookieModel } from './cookie.model';
import { map } from 'rxjs/operators';

import { ipAddress } from '../../../ip.conf';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent {
  ipAddy = ipAddress.ipFull; //localhost
  
  isLogin: boolean = false
  loginUsername: string | null = null
  loginPassword: string | null = null

  user = true;
  loading = false;
  submitted = false;
  form!: FormGroup;
  error = false;
  timer = false;

  constructor(
    private httpClient: HttpClient,
    private router: Router,
    private formBuilder: FormBuilder,
    private storageService: StorageService
  ){ }

  ngOnInit() {
    this.form = this.formBuilder.group({
        loginUsername: ['', Validators.required],
        loginPassword: ['', Validators.required]
        
    });
  } 

  get f() { 
    return this.form.controls; 
  }

  

  login(){
    this.submitted = true;
    type cookieData = 'session_token' | 'sessionToken' | 'expiresAt';
    this.httpClient.post<{session_token: string, sessionToken: string, expiresAt: string}>(this.ipAddy + '/signin1',  { //http://'+ this.ipAddy + ':8080/signin
      withCredentials: true,  
      username: this.loginUsername,
      password: this.loginPassword
    }).subscribe((response: {session_token: string, sessionToken: string, expiresAt: string}) => {
      console.log("login");
      
      const cookieResponse = response?.session_token + "=" + response?.sessionToken + "; Path=/; Expires=" + response?.expiresAt + ";"
      console.log(cookieResponse);

      document.cookie = cookieResponse;
      this.storageService.saveUser(this.loginUsername);

      this.timer = true;
      this.router.navigate(['profile'])
      this.loginUsername = null
      this.loginPassword = null
    });

    if(this.timer == false){
      this.error = true;
    }
  }

  goToPage(pageName:string){
    this.router.navigate([`${pageName}`]);
  }

  fieldTextType: boolean = false;
  show: boolean = false;

  toggleFieldTextType() {
      this.fieldTextType = !this.fieldTextType;
      this.show = !this.show;
  }

}
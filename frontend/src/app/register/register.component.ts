import { HttpClient } from '@angular/common/http';
import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ipAddress } from '../../../ip.conf';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css']
})
export class RegisterComponent {
  isLogin: boolean = false
  registerUsername: string | null = null
  registerPassword: string | null = null
  loginUsername: string | null = null
  loginPassword: string | null = null

  user = true;
  loading = false;
  submitted = false;
  form!: FormGroup;
  valid = false;
  error = false;

  constructor(
    private httpClient: HttpClient,
    private router: Router,
    private formBuilder: FormBuilder
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

  register(){
    this.submitted = true;
    console.log(this.registerUsername, this.registerPassword)
    this.httpClient.post(ipAddress.ipFull+'/signup1', {
      //username for now
      username: this.registerUsername,
      password: this.registerPassword
    }).pipe().subscribe((response) => {
      if(response){
      }
      //localStorage.setItem('token', response.jwt)
        this.router.navigate(['profile'])
        this.registerUsername = null
        this.registerPassword = null

        this.goToPage('signin');

        this.valid = true;
    })
    if(this.valid == false){
      this.error = true;
    }
  }

  goToPage(pageName:string){
    this.router.navigate([`${pageName}`]);
  }

}


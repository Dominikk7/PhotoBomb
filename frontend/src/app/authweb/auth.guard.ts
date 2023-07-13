import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import {
  ActivatedRouteSnapshot,
  CanActivate,
  Router,
} from "@angular/router";
import { ipAddress } from '../../../ip.conf';

@Injectable({
  providedIn: "root",
})
export class AuthGuard implements CanActivate {
  constructor(
    private router: Router,
    private httpClient: HttpClient
  ) {}

  canActivate(route: ActivatedRouteSnapshot) {
    console.log("Authgaurd");
    this.httpClient.get(ipAddress.ip+':4200/user-session') //'http://'+ipAddress.ip+':8080/user-session')
    .subscribe((res) => {
        console.log(res, 'AuthGuard')
    })
    if(localStorage.getItem('token')){
        return true
    }else{
        return false
    }
    
  }

}
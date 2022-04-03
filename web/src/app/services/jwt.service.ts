import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
@Injectable({
  providedIn: 'root'
})
export class JwtService {

  constructor(
    public http : HttpClient
  ) { }
  endpoint : string = 'http://localhost:8080'
  postUser = (user: string) =>{
    return this.http.post(`${this.endpoint}/signup`, user)
  }
  checkUser = (user: string): Observable<any> =>{
    return this.http.post(`${this.endpoint}/signin`, user, { withCredentials: true })
    //{ withCredentials: true } most imp thing otherwise token wont be saved at browser check config at server side 
    //a boolean value that indicates whether or not cross-site Access-Control requests should be made using credentials such as cookies, authorization headers or TLS client certificates
  }
  goHome = () =>{
    return this.http.get(`${this.endpoint}/home`, { withCredentials: true})
  }

  logout = () =>{
    console.log("...logging out")
    return this.http.get(`${this.endpoint}/logout`, { withCredentials: true})
  }
}

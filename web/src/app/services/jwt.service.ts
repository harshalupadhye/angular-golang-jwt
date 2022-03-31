import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
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
  checkUser = (user: string) =>{
    return this.http.post(`${this.endpoint}/signin`, user, { withCredentials: true })
    //{ withCredentials: true } most imp thing otherwise token wont be saved at browser check config at server side
  }
  goHome = () =>{
    return this.http.get(`${this.endpoint}/home`, { withCredentials: true })
  }
}

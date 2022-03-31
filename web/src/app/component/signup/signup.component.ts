import { Component, OnInit } from '@angular/core';
import { JwtService } from 'src/app/services/jwt.service';
import { Router } from '@angular/router';
@Component({
  selector: 'app-signup',
  templateUrl: './signup.component.html',
  styleUrls: ['./signup.component.scss'],
})
export class SignupComponent implements OnInit {
  username: string = '';
  password: string = '';
  isSignup: boolean = true
  constructor(public jwtService: JwtService, public router: Router) {}
  signup = () => {
    const user : Object = {
      username: this.username,
      password: this.password,
    };
    if(this.isSignup){
      this.jwtService.postUser(JSON.stringify(user)).subscribe((user) => {
        this.router.navigateByUrl('signin');
      });   
    } else {
      this.jwtService.checkUser(JSON.stringify(user)).subscribe((user: any) => {
        // document.cookie = `${user.Name}=${window.btoa(JSON.stringify(
        //   user.Value
        // ))};expires=${user.Expires}; domain=localhost; path=/`
        this.router.navigateByUrl('home');
      }); 
    }
   
  };

  ngOnInit(): void {
    this.isSignup = this.router.url == '/signup' ? true : false

  }
}

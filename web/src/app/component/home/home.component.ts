import { Component, OnInit } from '@angular/core';
import { JwtService } from 'src/app/services/jwt.service'; 
import { Router } from '@angular/router';
@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {
  username : any = null   
  constructor(public jwtService : JwtService, public router: Router) { }
  logout = () =>{
    let matchedCookie: string | any = document.cookie
      .match('(^|;)\\s*' + 'token' + '\\s*=\\s*([^;]+)')
      ?.pop()
      console.log(matchedCookie)
      document.cookie = `token=; Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT;`;
      this.router.navigateByUrl('signin')
  }
  ngOnInit(): void {
    this.jwtService.goHome().subscribe((user: any)=>{
      console.log('user',user)
      this.username = user  
    })
  }

}

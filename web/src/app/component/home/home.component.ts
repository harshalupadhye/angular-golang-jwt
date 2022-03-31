import { Component, OnInit } from '@angular/core';
import { JwtService } from 'src/app/services/jwt.service'; 

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {
  username : any = null   
  constructor(public jwtService : JwtService) { }

  ngOnInit(): void {
    this.jwtService.goHome().subscribe((user: any)=>{
      console.log('user',user)
      this.username = user  
    })
  }

}

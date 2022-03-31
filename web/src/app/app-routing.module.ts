import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { SignupComponent } from './component/signup/signup.component';
import { HomeComponent } from './component/home/home.component';
const routes: Routes = [
  {
    path:"signup", component:SignupComponent
  },
  {
    path:"signin", component:SignupComponent
  },
  {
    path:"home", component:HomeComponent
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }

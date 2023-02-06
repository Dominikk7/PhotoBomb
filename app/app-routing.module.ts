import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { PageComponent } from './page/page.component';

const routes: Routes = [
  {path: '', redirectTo: '/home', pathMatch: 'full'},
  {path: 'home', component: PageComponent, data: {
    page: 'home'
  }},
  {path: 'register', component: PageComponent, data: {
    page: 'register'
  }},
  {path: 'forget', component: PageComponent, data: {
    page: 'forget'
  }},
  {path: '**', redirectTo: '/home', pathMatch: 'full'}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})

export class AppRoutingModule { }

import { Routes } from '@angular/router';
import {LoginComponent} from "./home/login/login.component";
import {RegisterComponent} from "./home/register/register.component";

export const routes: Routes = [
	{
		title: 'Login',
		path: 'login',
		component: LoginComponent
	},
	{
		title: 'Register',
		path: 'register',
		component: RegisterComponent
	}
];

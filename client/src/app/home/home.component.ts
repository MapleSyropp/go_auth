import {Component, computed, Signal} from '@angular/core';
import {Router, RouterOutlet} from "@angular/router";
import {Button} from "primeng/button";
import {AuthService} from "../services/auth.service";

@Component({
  selector: 'app-home',
	imports: [
		RouterOutlet,
		Button
	],
  templateUrl: './home.component.html',
  styleUrl: './home.component.css'
})
export class HomeComponent {
	loginFlowStarted = computed(
		() => this.authService.getLoginFlowStatus()
	)

	constructor(private authService: AuthService, private router: Router) {
	}

	startLogin() {
		console.log("Llegue: ", this.loginFlowStarted());
		this.authService.setLoginFlowStatus(true);
		console.log("cambiado :", this.loginFlowStarted());
		this.router.navigate(['/login']);
	}
}

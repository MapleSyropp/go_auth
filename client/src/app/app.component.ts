import {Component, signal, WritableSignal} from '@angular/core';
import {HomeComponent} from './home/home.component';
import customTheme from "../assets/theme";
import {PrimeNG} from "primeng/config";
import {ToggleSwitch} from "primeng/toggleswitch";
import {FormsModule} from "@angular/forms";

@Component({
	selector: 'app-root',
	imports: [HomeComponent, ToggleSwitch, FormsModule],
	templateUrl: './app.component.html',
	styleUrl: './app.component.css'
})
export class AppComponent {
	darkMode: WritableSignal<boolean> = signal(true);

	constructor(private primeNgConfig: PrimeNG) {
		this.primeNgConfig.theme.set({
			preset: customTheme,
			options: {
				darkModeSelector: ".dark-mode",
			}
		});
	}

	toggleDarkMode() {
		const elem = document.querySelector('html');
		elem?.classList.toggle('dark-mode');
		this.darkMode.set(!this.darkMode());
	}
}

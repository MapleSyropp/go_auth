import { definePreset } from '@primeng/themes';
import Aura from '@primeng/themes/aura';

const customTheme = definePreset(Aura, {
	semantic: {
		primary: {
			50: '{amber.50}',
			100: '{amber.100}',
			200: '{amber.200}',
			300: '{amber.300}',
			400: '{amber.400}',
			500: '{amber.500}',
			600: '{amber.600}',
			700: '{amber.700}',
			800: '{amber.800}',
			900: '{amber.900}',
			950: '{amber.950}'
		},
		colorScheme: {
			light: {
				primary: {
					color: '{rose.600}',
					inverseColor: '#FFFFFF',
					hoverColor: '{rose.500}',
					activeColor: '{rose.400}',
				},
				highlight: {
					background: '{rose.600}',
					focusBackground: '{rose.300}',
					color: '#FFFFFF',
					focusColor: '#FFFFFF',
				},
				formField: {
					hoverBorderColor: '{amber.500}',
				}
			},
			dark: {
				primary: {
					color: '{amber.500}',
					inverseColor: '{amber.950}',
					hoverColor: '{amber.200}',
					activeColor: '{amber.300}',
				},
				highlight: {
					background: '{amber.600}',
					focusBackground: '{amber.800}',
					color: '{amber.950}',
					focusColor: '{amber.950}',
				},
				formField: {
					hoverBorderColor: '{amber.500}',
				}
			},
		}
	}
});

export default customTheme;

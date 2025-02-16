import i18next from "i18next";
import { initReactI18next } from "react-i18next";
import en from './locales/en.json';

const savedLng = localStorage.getItem('language') || 'en';

i18next.use(initReactI18next).init({
	resources: {
		en: { translation: en }
	},
	lng: savedLng,
	fallbackLng: 'en',
	interpolation: {escapeValue: false}
});

export default i18next;

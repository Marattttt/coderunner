import { useTranslation } from "react-i18next"

const Footer = () => {
	const {t} = useTranslation()

	return <footer className="
		flex flex-row justify-start items-center 
		left-0 bottom-0 p-4 md:p-8 w-full min-h-[64px] 
		font-heading text-xl 
		text-text-primary bg-bg-dark"
	>
		{t('footer.description')}
	</footer>
}

export default Footer

import { useTranslation } from "react-i18next";
import StackCard from './StackCard'

interface StackItem {
	title: string;
	items: {
		name: string;
		stack: string[]
	}[]
}

const About = () => {
	const {t} = useTranslation()

	const techStack = t('about.stack', {returnObjects: true}) as StackItem[]
	const cardContents: StackItem[] = techStack.map((st) => ({
		title: st.title,
		items: st.items
	}))


	return (<div className="flex flex-col justify-center align-top gap-4 w-full p-4 md:p-10 bg-bg-main" >
		<h2 className="text-text-primary text-center font-medium text-3xl">
			{t('about.title')}
		</h2>
		<div className="flex flex-col lg:flex-row gap-4 md:gap-8 w-full"> 
			{ cardContents.map((c) => <StackCard {...c} className="font-lg text-text-primary" key={c.title}/>) }
		</div>
	</div>)
}

export default About

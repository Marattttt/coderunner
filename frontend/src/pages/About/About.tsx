import React from "react";
import { useTranslation } from "react-i18next";

import GitHub from '../../assets/Github-White.svg?react'
import Telegram from '../../assets/Telegram.svg?react'

import StackCard from './StackCard'
import Contact from "./Contact";

interface StackItem {
	title: string;
	items: {
		name: string;
		stack: string[]
	}[]
}

const AboutHeader: React.FC<{ text: string }> = ({ text }) => {
	return <h2 className="text-text-primary text-center font-medium text-3xl">
		{text}
	</h2>
}

const About = () => {
	const { t } = useTranslation()

	const techStack = t('about.stack', { returnObjects: true }) as StackItem[]
	const cardContents: StackItem[] = techStack.map((st) => ({
		title: st.title,
		items: st.items
	}))

	return (<div className="flex flex-col justify-center align-top gap-4 w-full p-4 md:p-10 bg-bg-main" >
		<AboutHeader text={t('about.title')} />

		<div className="flex flex-col lg:flex-row gap-4 md:gap-8 w-full">
			{cardContents.map((c) => <StackCard {...c} className="font-lg text-text-primary" key={c.title} />)}
		</div>

		<AboutHeader text={t('about.contacts.title')} />
		<div className="flex flex-col md:flex-row gap-4 md:justify-center w-full text-text-primary">
			<Contact
				className="py-2 px-4 w-full rounded-xl border-2 border-bg-accent bg-bg-secondary"
				logo={<GitHub className="w-8 h-8" />}
				url="https://github.com/marattttt"
				text={t('about.contacts.github')}
			/>

			<Contact
				className="py-2 px-4 w-full rounded-xl border-2 border-bg-accent bg-bg-secondary"
				logo={<Telegram className="w-8 h-8" />}
				url="https://t.me/mtismyname"
				text={t('about.contacts.telegram')}
			/>
		</div>
	</div>)
}

export default About

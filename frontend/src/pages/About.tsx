import { useTranslation } from "react-i18next";

interface StackCardProps {
	title: string;
	items: {
		name: string;
		stack: string[];
	}[],
	className?: string;
}

const StackCard: React.FC<StackCardProps> = ({title,  items, className}) => {
	return (<section className={'flex flex-col w-full gap-2 p-4 sm:p-6 bg-bg-secondary ' + (className ?? '')}>
		<h3 className={'w-full text-center font-medium font-heading text-xl border-b-1 border-text-primary'}>
			{title}
		</h3>
		<div className="flex flex-row justify-between g-6 w-full ">
			{items.map((c) => (<div className={`flex flex-col w-full `} key={c.name}>
				<h4 className="font-lg font-medium text-xl pl-4">
					{c.name}
				</h4>
				<ol className="list-decimal pl-4 font-medium">
					{c.stack.map((i) => <li key={i}>{i}</li>)}
				</ol>
			</div>))}
		</div>
	</section>)
};

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
	const cardContents: StackCardProps[] = techStack.map((st) => ({
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

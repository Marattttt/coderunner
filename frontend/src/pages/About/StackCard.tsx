export interface StackCardProps {
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
		<div className="flex flex-row flex-wrap sm:flex-nowrap justify-between g-6 w-full ">
			{items.map((c) => (<div className={`flex flex-col w-full p-2`} key={c.name}>
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

export default StackCard

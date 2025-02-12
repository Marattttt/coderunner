import React from "react";

type ButtonProps = {
	onClick: () => void;
	children: React.ReactNode;
};

const ActionButton: React.FC<ButtonProps> = ({ onClick, children }) => {
	return (
		<button onClick={() => onClick()} 
			className="px-6 py-2 bg-bg-accent 
			font-bold text-xl text-text-primary rounded-xl 
			transition duration-200
			border-transparent
			border-1
			hover:border-text-primary"
		>
			{children}
		</button>
	);
};

export default ActionButton;

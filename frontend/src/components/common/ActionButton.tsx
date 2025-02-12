import React from "react";

type ButtonProps = {
	onClick: () => void;
	children: React.ReactNode;
};

const ActionButton: React.FC<ButtonProps> = ({ onClick, children }) => {
	return (
		<button onClick={() => onClick()} className="px-6 py-2 bg-bg-accent text-text-primary rounded-md hover:border-text-primary-2">
			{children}
		</button>
	);
};

export default ActionButton;

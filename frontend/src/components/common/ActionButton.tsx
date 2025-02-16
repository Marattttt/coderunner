import React from "react";

type ButtonProps = {
	onClick: () => void;
	children: React.ReactNode;
	disabled?: boolean
};

const ActionButton: React.FC<ButtonProps> = ({ onClick, children, disabled }) => {
	return (
		<button onClick={() => onClick()} 
			disabled={disabled}
			className={`
				px-6 py-2 bg-bg-accent 
				text-xl text-text-primary rounded-xl 
				transition duration-200
				border border-transparent
				hover:border-text-primary
				${disabled ? "opacity-50 cursor-not-allowed hover:border-transparent" : ""}
			`}
		>
			{children}
		</button>
	);
};

export default ActionButton;

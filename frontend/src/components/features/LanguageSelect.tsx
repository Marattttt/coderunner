import React, { useState } from "react";
import { Language } from "../../constants";

interface LanguageSelectProps {
	languages: Language[];
	onChange: (lang: Language) => void
}

const LanguageSelect: React.FC<LanguageSelectProps> = ({languages, onChange}) => {
	const [selected, setSelected] = useState(languages[0])
	const [isOpen, setIsOpen] = useState(false);

return(
    <div className="relative inline-block">
      {/* Button to toggle dropdown */}
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="px-4 py-2 bg-blue-500 text-white rounded-lg focus:outline-none"
      >
        {selected.name}
      </button>

      {/* Dropdown menu */}
      {isOpen && (
        <div className="absolute left-0 mt-2 w-48 bg-white shadow-lg rounded-lg border border-gray-200">
          {languages.map((lang) => (
            <button
              key={lang.id}
              className="block w-full px-4 py-2 text-left hover:bg-gray-100"
              onClick={() => {
                setSelected(lang);
                setIsOpen(false);
		if (selected !== lang) {
			onChange(lang)
		}
              }}
            >
              {lang.name}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}

export default LanguageSelect

import React, { useRef, useState } from "react";
import {
    Editor,
    EditorState,
    RichUtils,
} from "draft-js";
import Toolbar from "./Toolbar";
import "./RichTextEditor.css";


const myBlockStyleFn = (contentBlock) => {
    const type = contentBlock.getType();
    switch (type) {
        case "blockQuote":
            return "superFancyBlockquote";
        case "leftAlign":
            return "leftAlign";
        case "rightAlign":
            return "rightAlign";
        case "centerAlign":
            return "centerAlign";
        case "justifyAlign":
            return "justifyAlign";
        default:
            break;
    }
};

const styleMap = {
    CODE: {
        backgroundColor: "rgba(0, 0, 0, 0.05)",
        fontFamily: '"Inconsolata", "Menlo", "Consolas", monospace',
        fontSize: 16,
        padding: 2,
    },
    HIGHLIGHT: {
        backgroundColor: "#F7A5F7",
    },
    UPPERCASE: {
        textTransform: "uppercase",
    },
    LOWERCASE: {
        textTransform: "lowercase",
    },
    CODEBLOCK: {
        fontFamily: '"fira-code", "monospace"',
        fontSize: "inherit",
        background: "#ffeff0",
        fontStyle: "italic",
        lineHeight: 1.5,
        padding: "0.3rem 0.5rem",
        borderRadius: " 0.2rem",
    },
    SUPERSCRIPT: {
        verticalAlign: "super",
        fontSize: "80%",
    },
    SUBSCRIPT: {
        verticalAlign: "sub",
        fontSize: "80%",
    },
};

const RichTextEditor = ({ onChange }) => {
    const [editorState, setEditorState] = useState(
        EditorState.createEmpty()
    );
    const editor = useRef(null);

    const focusEditor = () => {
        editor.current.focus();
    };

    const handleKeyCommand = (command) => {
        const newState = RichUtils.handleKeyCommand(editorState, command);
        if (newState) {
            setEditorState(newState);
            return true;
        }
        return false;
    };

    return (
        <div className="rich-text-editor">
            <div className="editor-wrapper" onClick={focusEditor}>
                <Toolbar editorState={editorState} setEditorState={setEditorState} />
                <div className="editor-container">
                    <Editor
                        ref={editor}
                        placeholder=""
                        handleKeyCommand={handleKeyCommand}
                        editorState={editorState}
                        customStyleMap={styleMap}
                        blockStyleFn={myBlockStyleFn}
                        onChange={(editorState) => {
                            onChange(editorState);
                            setEditorState(editorState);
                        }}
                    />
                </div>
            </div>
        </div>
    );
};

export default RichTextEditor;

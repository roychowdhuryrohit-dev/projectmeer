import React, { useState, useEffect, useRef } from 'react';
import { LexicalComposer, PlainTextPlugin, $getRoot, $getSelection } from 'lexical';
import { $generateNodesFromDOM } from '@lexical/html';  // Optional:  If you will be pasting in content.
import {ContentEditable} from '@lexical/react/LexicalContentEditable';
import './Editor.css'
const initialConfig = {
  namespace: 'myEditor',
  theme: { /* Your custom theme */ },
  onError: (error) => console.error(error),
};

const Editor = () => {
  const [editorState, setEditorState] = useState(null);
  const editorRef = useRef(null);

  useEffect(() => {
    async function initializeEditor() {
      const startingEditorState = await editorRef.current.createEditorState();
      setEditorState(startingEditorState);
    }
    initializeEditor();
  }, []);

  function updateEditorState(editorState) {
    editorState.read(() => {
      const root = $getRoot();
      const selection = $getSelection();

      // ... Your custom update logic if needed ...

      setEditorState(editorState);
    });
  }

  function extractPlainText() {
    if (editorState) {
      editorState.read(() => {
        const textContent = $getRoot().getTextContent();
        console.log('Plain Text Content:', textContent);
      });
    }
  }

  return (
    <>
      <div className="editor-inner">
        
      {editorState && (
        <LexicalComposer initialConfig={initialConfig} editorState={editorState} onChange={updateEditorState} ref={editorRef}>
          <PlainTextPlugin contentEditable={<ContentEditable className="editor-input" />} />
        </LexicalComposer>
      )}
        </div>
        <button onClick={extractPlainText}>Get Plain Text</button>

    </>
  );
};

export default Editor;
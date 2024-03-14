import {createEditor, $createTextNode, $getRoot, $getSelection, $createParagraphNode} from 'lexical';
import {useDispatch, useSelector} from 'react-redux';
import {useParams} from 'react-router-dom';
import {useEffect, useState} from 'react';
import axios from 'axios';
import {LexicalComposer} from '@lexical/react/LexicalComposer';
import {PlainTextPlugin} from '@lexical/react/LexicalPlainTextPlugin';
import {ContentEditable} from '@lexical/react/LexicalContentEditable';
import {HistoryPlugin} from '@lexical/react/LexicalHistoryPlugin';
import {OnChangePlugin} from '@lexical/react/LexicalOnChangePlugin';
import {useLexicalComposerContext} from '@lexical/react/LexicalComposerContext';
import LexicalErrorBoundary from '@lexical/react/LexicalErrorBoundary';

import './Editor.css'



function MyOnChangePlugin({ onChange }) {
   
    const [editor] = useLexicalComposerContext();
       useEffect(() => {
           return editor.registerUpdateListener(({editorState}) => {
        onChange(editorState);
      });
    }, [editor, onChange]);
    return null;
}

function Editor() {
  const initialConfig = {
    namespace: 'MyEditor',
    onError: console.error,
  };

  const [editorState, setEditorState] = useState();
  const editor = createEditor(initialConfig);

  editor.update(() => {
    const root = $getRoot();
    const paragraphNode = $createParagraphNode();
    const textNode = $createTextNode();
    paragraphNode.append(textNode);
    root.append(paragraphNode);
  });
  
  function onChange(editorState) {
    // Call toJSON on the EditorState object, which produces a serialization safe string
    const editorStateJSON = editorState.toJSON();
    // However, we still have a JavaScript object, so we need to convert it to an actual string with JSON.stringify
    setEditorState(JSON.stringify(editorStateJSON));
  }
  return (
    <LexicalComposer initialConfig={initialConfig}>
      <div className="editor-inner">
      <PlainTextPlugin
        contentEditable={<ContentEditable />}
        placeholder={<div>Enter some text...</div>}
        ErrorBoundary={LexicalErrorBoundary}
      />
      <HistoryPlugin />
      <MyOnChangePlugin onChange={onChange}/>
      </div>
    </LexicalComposer>
  );
}

export default Editor;
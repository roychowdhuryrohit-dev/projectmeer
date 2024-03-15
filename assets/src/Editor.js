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
import Diff from './diff'

import './Editor.css'

const [previousEditorState, setPreviousEditorState] = useState(null);


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
  
  function updateEditorState(editorState) {
    editorState.update(() => {
      if (previousEditorState.editorState.editorState === undefined) {
        setPreviousEditorState(editorState);
        return;
      }
      oldText = previousEditorState.$getRoot().getTextContent();
      newText = editorState.$getRoot().getTextContent();
    
      root.clear();
      const paragraph = $createParagraphNode();
      paragraph.append($createTextNode(newText));
      root.append(paragraph);
      
      setEditorState(editorState);
    });
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
      <MyOnChangePlugin onChange={updateEditorState}/>
      </div>

    </LexicalComposer>
  );
}

export default Editor;
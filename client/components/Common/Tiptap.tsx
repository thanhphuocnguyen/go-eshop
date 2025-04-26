import React, { useEffect } from 'react';
import './styles.css';
import { Color } from '@tiptap/extension-color';
import ListItem from '@tiptap/extension-list-item';
import TextStyle from '@tiptap/extension-text-style';
import { Editor, EditorContent, useEditor } from '@tiptap/react';
import StarterKit from '@tiptap/starter-kit';
import {
  Button,
  Combobox,
  ComboboxButton,
  ComboboxInput,
  ComboboxOption,
  ComboboxOptions,
} from '@headlessui/react';
import {
  CheckIcon,
  ChevronDownIcon,
  CodeBracketIcon,
  BoltIcon,
  ListBulletIcon,
  QueueListIcon,
  ArrowUturnLeftIcon,
  ArrowUturnRightIcon,
  MinusIcon,
} from '@heroicons/react/24/outline';
import clsx from 'clsx';
import {
  Control,
  Controller,
  FieldValues,
  Path,
  PathValue,
} from 'react-hook-form';

const extensions = [
  Color.configure({ types: [TextStyle.name, ListItem.name] }),
  TextStyle.configure({}),
  StarterKit.configure({
    bulletList: {
      keepMarks: true,
      keepAttributes: false, // TODO : Making this as `false` becase marks are not preserved when I try to preserve attrs, awaiting a bit of help
    },
    orderedList: {
      keepMarks: true,
      keepAttributes: false, // TODO : Making this as `false` becase marks are not preserved when I try to preserve attrs, awaiting a bit of help
    },
  }),
];

interface TiptapEditorProps {
  content: string;
  onChange: (content: string) => void;
}
export const TiptapEditor: React.FC<TiptapEditorProps> = ({
  content,
  onChange,
}) => {
  const editor = useEditor({
    extensions,
    content,
    onUpdate: ({ editor }) => {
      const html = editor.getHTML();
      onChange(html);
    },
  });

  return (
    <div className='relative tiptap-editor'>
      <MenuBar editor={editor} />
      <EditorContent className='bg-white' content={content} editor={editor} />
    </div>
  );
};

interface TiptapControllerProps<T extends FieldValues> {
  name: Path<T>;
  control: Control<T>;
  defaultValue?: PathValue<T, Path<T>>;
  label?: string;
  error?: string;
}

export const TiptapController = <T extends FieldValues>({
  name,
  control,
  defaultValue = '' as PathValue<T, Path<T>>,
  label,
  error,
}: TiptapControllerProps<T>) => {
  return (
    <div className='w-full'>
      {label && (
        <label className='text-sm/6 text-gray-500 font-semibold'>{label}</label>
      )}
      <Controller
        name={name}
        control={control}
        defaultValue={defaultValue}
        render={({ field: { onChange, value } }) => (
          <div className={clsx('mt-1', error && 'border-red-500')}>
            <TiptapEditor
              content={value || ''}
              onChange={(html) => onChange(html)}
            />
          </div>
        )}
      />
      {error && <p className={clsx('text-sm mt-2', 'text-red-500')}>{error}</p>}
    </div>
  );
};

const options = [
  { id: 'paragraph', name: 'Paragraph' },
  { id: 'h1', name: 'Heading 1' },
  { id: 'h2', name: 'Heading 2' },
  { id: 'h3', name: 'Heading 3' },
  { id: 'h4', name: 'Heading 4' },
  { id: 'h5', name: 'Heading 5' },
];

const MenuBar = ({ editor }: { editor: Editor | null }) => {
  const [selected, setSelected] = React.useState(options[0]);

  useEffect(() => {
    if (selected && editor) {
      switch (selected.id) {
        case 'paragraph':
          editor?.chain().focus().setParagraph().run();
          break;
        case 'h1':
          editor?.chain().focus().toggleHeading({ level: 1 }).run();
          break;
        case 'h2':
          editor?.chain().focus().toggleHeading({ level: 2 }).run();
          break;
        case 'h3':
          editor?.chain().focus().toggleHeading({ level: 3 }).run();
          break;
        case 'h4':
          editor?.chain().focus().toggleHeading({ level: 4 }).run();
          break;
        case 'h5':
          editor?.chain().focus().toggleHeading({ level: 5 }).run();
          break;
        default:
          break;
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selected]);

  if (!editor) {
    return null;
  }

  return (
    <div className='quill-toolbar p-1 bg-gray-50 border border-gray-300 rounded-t-md'>
      <div className='flex flex-wrap items-center gap-0.5'>
        {/* Text format group */}
        <div className='toolbar-group flex border-r border-gray-300 mr-1 pr-1'>
          <Button
            onClick={() => editor.chain().focus().toggleBold().run()}
            disabled={!editor.can().chain().focus().toggleBold().run()}
            className={clsx(
              'toolbar-btn p-1 hover:bg-gray-200 rounded',
              editor.isActive('bold') ? 'bg-gray-200 font-bold' : ''
            )}
            title='Bold'
          >
            <span className='font-bold'>B</span>
          </Button>
          <Button
            onClick={() => editor.chain().focus().toggleItalic().run()}
            disabled={!editor.can().chain().focus().toggleItalic().run()}
            className={clsx(
              'toolbar-btn p-1 hover:bg-gray-200 rounded',
              editor.isActive('italic') ? 'bg-gray-200 italic' : ''
            )}
            title='Italic'
          >
            <span className='italic'>I</span>
          </Button>
          <Button
            onClick={() => editor.chain().focus().toggleStrike().run()}
            disabled={!editor.can().chain().focus().toggleStrike().run()}
            className={clsx(
              'toolbar-btn p-1 hover:bg-gray-200 rounded',
              editor.isActive('strike') ? 'bg-gray-200 line-through' : ''
            )}
            title='Strikethrough'
          >
            <span className='line-through'>S</span>
          </Button>
          <Button
            onClick={() => editor.chain().focus().toggleCode().run()}
            disabled={!editor.can().chain().focus().toggleCode().run()}
            className={clsx(
              'toolbar-btn p-1 hover:bg-gray-200 rounded',
              editor.isActive('code') ? 'bg-gray-200' : ''
            )}
            title='Code'
          >
            <CodeBracketIcon className='size-4' />
          </Button>
        </div>

        {/* Headings dropdown */}
        <div className='toolbar-group flex border-r border-gray-300 mr-1 pr-1'>
          <Combobox
            value={selected}
            onChange={(value) => value && setSelected(value)}
          >
            <div className='relative w-40'>
              <ComboboxInput
                disabled
                className={clsx(
                  'w-32 rounded bg-transparent px-2 py-1 text-sm text-gray-700 border border-gray-300',
                  'focus:outline-none focus:ring-1 focus:ring-sky-400'
                )}
                displayValue={(option) => option?.name ?? ''}
                // readOnly
              />
              <ComboboxButton className='absolute inset-y-0 right-0 px-2'>
                <ChevronDownIcon className='size-4 text-gray-500' />
              </ComboboxButton>
            </div>

            <ComboboxOptions
              anchor='bottom'
              className={clsx(
                'w-40 rounded border border-gray-300 bg-white p-1 shadow-lg z-10'
              )}
            >
              {options.map((option) => (
                <ComboboxOption
                  key={option.id}
                  value={option}
                  className='group flex cursor-pointer items-center gap-2 rounded py-1 px-2 hover:bg-gray-100'
                >
                  <CheckIcon className='invisible size-4 text-sky-500 group-data-[selected]:visible' />
                  <span className='text-sm text-gray-700'>{option.name}</span>
                </ComboboxOption>
              ))}
            </ComboboxOptions>
          </Combobox>
        </div>

        {/* Lists group */}
        <div className='toolbar-group flex border-r border-gray-300 mr-1 pr-1'>
          <Button
            onClick={() => editor.chain().focus().toggleBulletList().run()}
            className={clsx(
              'toolbar-btn p-1 hover:bg-gray-200 rounded',
              editor.isActive('bulletList') ? 'bg-gray-200' : ''
            )}
            title='Bullet List'
          >
            <ListBulletIcon className='size-4' />
          </Button>
          <Button
            onClick={() => editor.chain().focus().toggleOrderedList().run()}
            className={clsx(
              'toolbar-btn p-1 hover:bg-gray-200 rounded',
              editor.isActive('orderedList') ? 'bg-gray-200' : ''
            )}
            title='Ordered List'
          >
            <QueueListIcon className='size-4' />
          </Button>
        </div>

        {/* Block elements */}
        <div className='toolbar-group flex border-r border-gray-300 mr-1 pr-1'>
          <Button
            onClick={() => editor.chain().focus().toggleBlockquote().run()}
            className={clsx(
              'toolbar-btn p-1 hover:bg-gray-200 rounded flex items-center',
              editor.isActive('blockquote') ? 'bg-gray-200' : ''
            )}
            title='Blockquote'
          >
            <span className='text-lg'>&quot;</span>
          </Button>
          <Button
            onClick={() => editor.chain().focus().setHorizontalRule().run()}
            className='toolbar-btn p-1 hover:bg-gray-200 rounded'
            title='Horizontal Rule'
          >
            <MinusIcon className='size-4' />
          </Button>
        </div>

        {/* History controls */}
        <div className='toolbar-group flex'>
          <Button
            onClick={() => editor.chain().focus().undo().run()}
            disabled={!editor.can().chain().focus().undo().run()}
            className={clsx(
              'toolbar-btn p-1 hover:bg-gray-200 rounded',
              !editor.can().chain().focus().undo().run() &&
                'opacity-50 cursor-not-allowed'
            )}
            title='Undo'
          >
            <ArrowUturnLeftIcon className='size-4' />
          </Button>
          <Button
            onClick={() => editor.chain().focus().redo().run()}
            disabled={!editor.can().chain().focus().redo().run()}
            className={clsx(
              'toolbar-btn p-1 hover:bg-gray-200 rounded',
              !editor.can().chain().focus().redo().run() &&
                'opacity-50 cursor-not-allowed'
            )}
            title='Redo'
          >
            <ArrowUturnRightIcon className='size-4' />
          </Button>
        </div>

        {/* Clear formatting */}
        <div className='toolbar-group flex ml-auto'>
          <Button
            onClick={() => {
              editor.chain().focus().unsetAllMarks().run();
              editor.chain().focus().clearNodes().run();
            }}
            className='toolbar-btn p-1 hover:bg-gray-200 rounded text-xs'
            title='Clear formatting'
          >
            <BoltIcon className='size-4' /> Clear format
          </Button>
        </div>
      </div>
    </div>
  );
};

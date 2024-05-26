import { useRef } from 'react';
import { QuestionMarkCircleIcon } from '@heroicons/react/24/solid';

interface Props {
  hint: string;
}

export function HintButton({ hint }: Props): JSX.Element {
  const modalRef = useRef<HTMLDialogElement>(null);

  function onClick(e: React.MouseEvent<HTMLButtonElement>) {
    e.preventDefault();
    modalRef.current?.showModal();
  }

  function onClose(e: React.MouseEvent<HTMLButtonElement>) {
    e.preventDefault();
    modalRef.current?.close();
  }

  return (
    <>
      <button type="button" className="btn btn-ghost btn-circle" onClick={onClick}>
        <QuestionMarkCircleIcon className="w-8 h-8" />
      </button>
      <dialog ref={modalRef} className="modal modal-bottom sm:modal-middle">
        <div className="modal-box">
          <h3 className="font-bold text-lg">Hint</h3>
          <p className="py-4">{hint}</p>
          <div className="modal-action">
            <button onClick={onClose} className="btn">
              Close
            </button>
          </div>
        </div>
      </dialog>
    </>
  );
}

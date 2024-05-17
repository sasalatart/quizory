import { useRef } from 'react';

interface Props {
  hint: string;
}

export function HintButton({ hint }: Props): JSX.Element {
  const modalRef = useRef<HTMLDialogElement>(null);

  function onClick() {
    modalRef.current?.showModal();
  }

  return (
    <>
      <button className="btn btn-secondary" onClick={onClick}>
        Hint
      </button>
      <dialog ref={modalRef} className="modal modal-bottom sm:modal-middle">
        <div className="modal-box">
          <h3 className="font-bold text-lg">Hint</h3>
          <p className="py-4">{hint}</p>
          <div className="modal-action">
            <form method="dialog">
              {/* if there is a button in form, it will close the modal */}
              <button className="btn">Close</button>
            </form>
          </div>
        </div>
      </dialog>
    </>
  );
}

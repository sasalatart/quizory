import { ReactNode, RefObject } from 'react';
import { clsx } from 'clsx';
import { type AnswersLogItem } from '@/generated/api';
import { QuestionBadges } from '@/layout';

interface Props {
  logItem: AnswersLogItem | undefined;
  onClose: () => unknown;
  modalRef: RefObject<HTMLDialogElement>;
}

export function LogItemModal({ logItem, onClose, modalRef }: Props): JSX.Element | null {
  if (!logItem) {
    return null;
  }

  const { question } = logItem;
  return (
    <dialog ref={modalRef} className="modal modal-bottom sm:modal-middle">
      <div className="modal-box">
        <div className="mb-4">
          <QuestionBadges topic={question.topic} difficulty={question.difficulty} />
        </div>

        <Collapsible tabIndex={0} title="Question">
          <h3 className="font-bold text-lg">{question.question}</h3>
          <div className="mt-4">
            {question.choices.map(({ id, choice, isCorrect }) => {
              const isChecked = id === logItem.choiceId;
              return (
                <div key={id} className="form-control">
                  <label className="flex justify-start cursor-pointer label">
                    <input
                      name="choice"
                      type="radio"
                      value={id}
                      className="radio"
                      disabled
                      checked={isChecked}
                    />
                    <span
                      className={clsx(
                        'label-text',
                        'ml-4',
                        isCorrect && 'text-success',
                        isChecked && !isCorrect && 'text-error',
                      )}
                    >
                      {choice}
                    </span>
                  </label>
                </div>
              );
            })}
          </div>
        </Collapsible>

        <Collapsible tabIndex={1} title="Hint">
          <p>{question.hint}</p>
        </Collapsible>

        <Collapsible tabIndex={2} title="More info">
          {question.moreInfo.split('\n').map((line, index) => (
            <p key={index} className="my-4">
              {line}
            </p>
          ))}
        </Collapsible>

        <div className="modal-action">
          <button onClick={onClose} className="btn btn-block">
            Close
          </button>
        </div>
      </div>
    </dialog>
  );
}

interface CollapsibleProps {
  tabIndex: number;
  title: string;
  children: ReactNode;
}

function Collapsible({ title, children, tabIndex }: CollapsibleProps): JSX.Element {
  return (
    <div tabIndex={tabIndex} className="collapse">
      <p className="collapse-title text-xl font-medium link">{title}</p>
      <div className="collapse-content">{children}</div>
    </div>
  );
}

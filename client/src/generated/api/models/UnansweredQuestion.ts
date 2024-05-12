/* tslint:disable */
/* eslint-disable */
/**
 * AI Generated Questions
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 1.0.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { mapValues } from '../runtime';
import type { Difficulty } from './Difficulty';
import {
    DifficultyFromJSON,
    DifficultyFromJSONTyped,
    DifficultyToJSON,
} from './Difficulty';
import type { UnansweredChoice } from './UnansweredChoice';
import {
    UnansweredChoiceFromJSON,
    UnansweredChoiceFromJSONTyped,
    UnansweredChoiceToJSON,
} from './UnansweredChoice';

/**
 * 
 * @export
 * @interface UnansweredQuestion
 */
export interface UnansweredQuestion {
    /**
     * 
     * @type {string}
     * @memberof UnansweredQuestion
     */
    id: string;
    /**
     * 
     * @type {string}
     * @memberof UnansweredQuestion
     */
    topic: string;
    /**
     * 
     * @type {string}
     * @memberof UnansweredQuestion
     */
    question: string;
    /**
     * 
     * @type {string}
     * @memberof UnansweredQuestion
     */
    hint: string;
    /**
     * 
     * @type {Difficulty}
     * @memberof UnansweredQuestion
     */
    difficulty: Difficulty;
    /**
     * 
     * @type {Array<UnansweredChoice>}
     * @memberof UnansweredQuestion
     */
    choices: Array<UnansweredChoice>;
}

/**
 * Check if a given object implements the UnansweredQuestion interface.
 */
export function instanceOfUnansweredQuestion(value: object): boolean {
    if (!('id' in value)) return false;
    if (!('topic' in value)) return false;
    if (!('question' in value)) return false;
    if (!('hint' in value)) return false;
    if (!('difficulty' in value)) return false;
    if (!('choices' in value)) return false;
    return true;
}

export function UnansweredQuestionFromJSON(json: any): UnansweredQuestion {
    return UnansweredQuestionFromJSONTyped(json, false);
}

export function UnansweredQuestionFromJSONTyped(json: any, ignoreDiscriminator: boolean): UnansweredQuestion {
    if (json == null) {
        return json;
    }
    return {
        
        'id': json['id'],
        'topic': json['topic'],
        'question': json['question'],
        'hint': json['hint'],
        'difficulty': DifficultyFromJSON(json['difficulty']),
        'choices': ((json['choices'] as Array<any>).map(UnansweredChoiceFromJSON)),
    };
}

export function UnansweredQuestionToJSON(value?: UnansweredQuestion | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'id': value['id'],
        'topic': value['topic'],
        'question': value['question'],
        'hint': value['hint'],
        'difficulty': DifficultyToJSON(value['difficulty']),
        'choices': ((value['choices'] as Array<any>).map(UnansweredChoiceToJSON)),
    };
}


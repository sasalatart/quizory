/* tslint:disable */
/* eslint-disable */
/**
 * Quizory
 * LLM-Generated history questions to test your knowledge.
 *
 * The version of the OpenAPI document: 1.0.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { mapValues } from '../runtime';
/**
 * 
 * @export
 * @interface UnansweredChoice
 */
export interface UnansweredChoice {
    /**
     * 
     * @type {string}
     * @memberof UnansweredChoice
     */
    id: string;
    /**
     * 
     * @type {string}
     * @memberof UnansweredChoice
     */
    choice: string;
}

/**
 * Check if a given object implements the UnansweredChoice interface.
 */
export function instanceOfUnansweredChoice(value: object): boolean {
    if (!('id' in value)) return false;
    if (!('choice' in value)) return false;
    return true;
}

export function UnansweredChoiceFromJSON(json: any): UnansweredChoice {
    return UnansweredChoiceFromJSONTyped(json, false);
}

export function UnansweredChoiceFromJSONTyped(json: any, ignoreDiscriminator: boolean): UnansweredChoice {
    if (json == null) {
        return json;
    }
    return {
        
        'id': json['id'],
        'choice': json['choice'],
    };
}

export function UnansweredChoiceToJSON(value?: UnansweredChoice | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'id': value['id'],
        'choice': value['choice'],
    };
}


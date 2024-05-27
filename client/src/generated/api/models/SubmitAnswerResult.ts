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
 * @interface SubmitAnswerResult
 */
export interface SubmitAnswerResult {
    /**
     * 
     * @type {string}
     * @memberof SubmitAnswerResult
     */
    id: string;
    /**
     * 
     * @type {string}
     * @memberof SubmitAnswerResult
     */
    correctChoiceId: string;
    /**
     * 
     * @type {string}
     * @memberof SubmitAnswerResult
     */
    moreInfo: string;
}

/**
 * Check if a given object implements the SubmitAnswerResult interface.
 */
export function instanceOfSubmitAnswerResult(value: object): boolean {
    if (!('id' in value)) return false;
    if (!('correctChoiceId' in value)) return false;
    if (!('moreInfo' in value)) return false;
    return true;
}

export function SubmitAnswerResultFromJSON(json: any): SubmitAnswerResult {
    return SubmitAnswerResultFromJSONTyped(json, false);
}

export function SubmitAnswerResultFromJSONTyped(json: any, ignoreDiscriminator: boolean): SubmitAnswerResult {
    if (json == null) {
        return json;
    }
    return {
        
        'id': json['id'],
        'correctChoiceId': json['correctChoiceId'],
        'moreInfo': json['moreInfo'],
    };
}

export function SubmitAnswerResultToJSON(value?: SubmitAnswerResult | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'id': value['id'],
        'correctChoiceId': value['correctChoiceId'],
        'moreInfo': value['moreInfo'],
    };
}


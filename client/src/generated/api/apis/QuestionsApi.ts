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


import * as runtime from '../runtime';
import type {
  RemainingTopic,
  UnansweredQuestion,
} from '../models/index';
import {
    RemainingTopicFromJSON,
    RemainingTopicToJSON,
    UnansweredQuestionFromJSON,
    UnansweredQuestionToJSON,
} from '../models/index';

export interface GetNextQuestionRequest {
    topic: string;
}

/**
 * 
 */
export class QuestionsApi extends runtime.BaseAPI {

    /**
     * Returns the next question that a user should answer for the specified topic.
     */
    async getNextQuestionRaw(requestParameters: GetNextQuestionRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<UnansweredQuestion>> {
        if (requestParameters['topic'] == null) {
            throw new runtime.RequiredError(
                'topic',
                'Required parameter "topic" was null or undefined when calling getNextQuestion().'
            );
        }

        const queryParameters: any = {};

        if (requestParameters['topic'] != null) {
            queryParameters['topic'] = requestParameters['topic'];
        }

        const headerParameters: runtime.HTTPHeaders = {};

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("BearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/questions/next`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => UnansweredQuestionFromJSON(jsonValue));
    }

    /**
     * Returns the next question that a user should answer for the specified topic.
     */
    async getNextQuestion(requestParameters: GetNextQuestionRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<UnansweredQuestion | null | undefined > {
        const response = await this.getNextQuestionRaw(requestParameters, initOverrides);
        switch (response.raw.status) {
            case 200:
                return await response.value();
            case 204:
                return null;
            default:
                return await response.value();
        }
    }

    /**
     * Returns the list of topics with questions still unanswered by the user making the request. Each of these topics comes with the actual amount of questions left to answer. 
     */
    async getRemainingTopicsRaw(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Array<RemainingTopic>>> {
        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("BearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/questions/remaining-topics`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => jsonValue.map(RemainingTopicFromJSON));
    }

    /**
     * Returns the list of topics with questions still unanswered by the user making the request. Each of these topics comes with the actual amount of questions left to answer. 
     */
    async getRemainingTopics(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Array<RemainingTopic>> {
        const response = await this.getRemainingTopicsRaw(initOverrides);
        return await response.value();
    }

}

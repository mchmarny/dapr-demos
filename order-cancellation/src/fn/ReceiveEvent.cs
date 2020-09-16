// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

namespace DaprAzFn.Sample
{
    using System;
    using System.IO;
    using System.Threading.Tasks;
    using Microsoft.AspNetCore.Mvc;
    using CloudNative.CloudEvents;
    using Microsoft.Azure.WebJobs;
    using Dapr.AzureFunctions.Extension;
    using Microsoft.Azure.WebJobs.Extensions.Http;
    using Microsoft.AspNetCore.Http;
    using Microsoft.Extensions.Logging;
    using Newtonsoft.Json.Linq;

    public static class ReceiveTopicMessage
    {
        /// <summary>
        /// Subscribes to topic and saves it to state store 
        /// </summary>
        [FunctionName("ReceiveTopicMessage")]
        public static async Task<IActionResult> Run(
            [DaprTopicTrigger("%PubSubName%", Topic = "%TopicName%")] CloudEvent cloudEvent,
            [DaprState("%StateStore%")] IAsyncCollector<DaprStateRecord> state, ILogger log)
        {
            // Get data from CloudEvent
            log.LogInformation($"Received message: {cloudEvent.Data}.");
            var cancellationData = cloudEvent.Data as JToken;

            // TODO: Implement cancellation validaiton logic.  

            // Parse state ID and Save
            var keyname = Environment.GetEnvironmentVariable("StateKey") ?? "id";
            var key = cancellationData.Value<string>(keyname);
            var stateRec = new DaprStateRecord(key, cancellationData);
            await state.AddAsync(stateRec);
            return new OkResult();
        }
    }
}

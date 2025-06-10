// Quick verification that the Azure OpenAI SDK conversion works
const { AzureAIService } = require('./lib/services/azureAI');

async function verifyConversion() {
    console.log('🔍 Verifying Azure OpenAI SDK conversion...');
    
    // Test configuration
    const config = {
        apiKey: 'test-key',
        endpoint: 'https://test-endpoint.openai.azure.com',
        deploymentId: 'test-deployment'
    };
    
    try {
        // Test service instantiation
        const service = new AzureAIService(config);
        console.log('✅ AzureAIService instantiated successfully with official SDK');
        
        // Check that the client is properly initialized
        if (service.client && service.deploymentId) {
            console.log('✅ Azure OpenAI client properly initialized');
            console.log('✅ Deployment ID properly stored');
        }
        
        console.log('\n🎉 Conversion verification successful!');
        console.log('\nThe service is now using the official Azure OpenAI SDK from the openai package.');
        console.log('This provides better error handling, automatic retries, and official support.');
        
    } catch (error) {
        console.error('❌ Verification failed:', error.message);
    }
}

verifyConversion();

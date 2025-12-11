import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

export default function Consent() {
  return (
    <div className="container max-w-3xl py-16">
      <Card>
        <CardHeader>
          <CardTitle className="text-2xl text-center font-bold">
            CONSENT TO THE PROCESSING OF PERSONAL DATA
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4 text-sm leading-relaxed text-muted-foreground">
          <p>
            Acting freely, by my own will and in my own interest, and also confirming my legal
            capacity, the individual gives their consent to Krivonosov Konstantin Yurievich,
            registered at the address: 350912, Krasnodar Krai, Krasnodar city, 412 Evdokiya
            Bershanskaya Street, Building 179, 10th floor (Hereinafter - the Operator) to process
            their personal data under the following conditions:
          </p>
          <ul className="list-disc pl-5 space-y-2">
            <li>
              This Consent is given for the processing of personal data, both without the use of
              automation tools and with their use.
            </li>
            <li>
              Consent is given for the processing of the following personal data of mine: data
              available from the social network profile used for authorization on the site; email
              address; nickname.
            </li>
            <li>
              Purpose of personal data processing: registration of the user on the site/in the
              application.
            </li>
            <li>
              In the course of processing, the following actions will be performed with personal
              data: collection, systematization; storage; use; retrieval; blocking; destruction;
              recording; deletion; accumulation; updating; modification.
            </li>
            <li>
              Personal data is processed until the user deletes their personal account on the site.
            </li>
            <li>
              Consent may be revoked by the personal data subject or their representative by sending
              an application to the email address:{' '}
              <a href="mailto:koskriv2006@gmail.com" className="text-primary hover:underline">
                koskriv2006@gmail.com
              </a>
              .
            </li>
          </ul>
        </CardContent>
      </Card>
    </div>
  );
}

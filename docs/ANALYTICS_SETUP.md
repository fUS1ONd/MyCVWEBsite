# Analytics Setup (GTM + Yandex Metrica)

Google Tag Manager (GTM) is integrated directly into `index.html` in this project. For correct Single Page Application (SPA) operation, a virtual page view tracking mechanism has been implemented.

## 1. Environment Variables Setup

In the `frontend/.env` file (or deployment settings), specify the GTM container ID:

```bash
VITE_GTM_ID=GTM-XXXXXXX
```

> **Important:** If environment variables are not automatically substituted (for example, during local development without a special plugin), you need to manually replace `%VITE_GTM_ID%` in `frontend/index.html` with your container ID.

## 2. Yandex.Metrica Setup Inside GTM

Since the React application works as an SPA, the standard Metrica code will not track page transitions. Set up tags as follows:

### Step 1: Creating a Trigger

1. Go to the **Triggers** section.
2. Create a new trigger.
3. Trigger type: **Custom Event**.
4. Event Name: `page_view`.
   * *This event is sent by our code (`frontend/src/lib/gtm.ts`) on every route change.*
5. Save the trigger with a name like "Custom Page View".

### Step 2: Creating a Variable (Optional)

1. If you need to pass the exact URL, go to **Variables**.
2. Create a variable of type **Data Layer Variable**.
3. Data Layer Variable Name: `page_path`.

### Step 3: Creating the Metrica Tag

1. Go to the **Tags** section.
2. Create a new tag of type **Custom HTML**.
3. Insert the Yandex.Metrica counter code.
4. **Important:** Make sure that automatic tracking is disabled in the counter code (if it conflicts) or configured correctly. For SPA, it's more reliable to call the `hit` method explicitly.
5. Inside the script, add `hit` processing when the tag fires:

   ```html
   <script>
      ym(YOUR_COUNTER_ID, 'hit', {{page_path}}); // Use variable from Step 2 or window.location.href
   </script>
   ```
6. In the **Triggering** section, select the trigger created in Step 1 ("Custom Page View").

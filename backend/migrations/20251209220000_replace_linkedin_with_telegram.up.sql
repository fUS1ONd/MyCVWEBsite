UPDATE profile_info
SET contacts = (contacts - 'linkedin') || jsonb_build_object('telegram', 'https://t.me/fUS1ONd');

UPDATE profile_info
SET contacts = (contacts - 'telegram') || jsonb_build_object('linkedin', '');

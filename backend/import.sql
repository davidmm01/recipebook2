-- Auto-generated SQL import script
-- Generated at: 2025-11-04T18:03:24+11:00

BEGIN TRANSACTION;

-- Recipe: Boulevardier
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('9dbb6523-bdd2-40ac-acd8-9a0cc50fb711', 'Boulevardier', '', 'drink', 'string', '- 1 part bourbon
- 1 part campari
- 1 part vermouth rosso
- orange slice or peel for garnish', '- Combine.
- Garnish with orange slice.
- Serve on the rocks.', '', '', 'dave', '2025-09-28 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('bourbon');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '9dbb6523-bdd2-40ac-acd8-9a0cc50fb711', id FROM tags WHERE name = 'bourbon';
INSERT OR IGNORE INTO tags (name) VALUES ('campari');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '9dbb6523-bdd2-40ac-acd8-9a0cc50fb711', id FROM tags WHERE name = 'campari';
INSERT OR IGNORE INTO tags (name) VALUES ('vermouth');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '9dbb6523-bdd2-40ac-acd8-9a0cc50fb711', id FROM tags WHERE name = 'vermouth';

-- Recipe: Gin Jam Fizz
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('f15062be-0675-4897-86bc-20c1193ce9b8', 'Gin Jam Fizz', '', 'drink', 'string', '- 45ml pink gin
- 2 tsp raspberry jam
- 1/2 cup ice
- 15ml lemon juice', '- Shake
- Top up the (wine) glass with soda water to liking.', '### Notes
- Could easily add 1 tsp of jam or sugar syrup for more sweetness
- could use 60ml gin if you want it strong', '', 'megasaur', '2023-12-27 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('gin');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'f15062be-0675-4897-86bc-20c1193ce9b8', id FROM tags WHERE name = 'gin';

-- Recipe: Margarita
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('531d6c7a-9a19-4362-9195-ba2e1d226cac', 'Margarita', '', 'drink', 'string', '- juice of 1 lime
- 45ml tequila (blanco or reposado)
- 15ml cointreu
- 15ml agave syrup
- pinch of salt', '- Combine all ingredients.
- Shake
- Strain and serve.
- Optionally serve ice and or a splash of soda water.', '### Notes
- Make no more than 2 margaritas into a standard shaker at a time, else you will not be able to add enough ice to the shaker.', '', 'dave', '2025-09-28 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('tequila');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '531d6c7a-9a19-4362-9195-ba2e1d226cac', id FROM tags WHERE name = 'tequila';
INSERT OR IGNORE INTO tags (name) VALUES ('lime');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '531d6c7a-9a19-4362-9195-ba2e1d226cac', id FROM tags WHERE name = 'lime';
INSERT OR IGNORE INTO tags (name) VALUES ('salt');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '531d6c7a-9a19-4362-9195-ba2e1d226cac', id FROM tags WHERE name = 'salt';

-- Recipe: Negroni
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('c8e95760-1f47-4f13-ab38-4f2c766658b8', 'Negroni', '', 'drink', 'string', '- 1 part gin
- 1 part campari
- 1 part vermouth rosso
- orange slice for garnish', '- Combine.
- Garnish with orange slice.
- Serve on the rocks.', '', '', 'dave', '2025-09-28 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('gin');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'c8e95760-1f47-4f13-ab38-4f2c766658b8', id FROM tags WHERE name = 'gin';
INSERT OR IGNORE INTO tags (name) VALUES ('campari');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'c8e95760-1f47-4f13-ab38-4f2c766658b8', id FROM tags WHERE name = 'campari';
INSERT OR IGNORE INTO tags (name) VALUES ('vermouth');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'c8e95760-1f47-4f13-ab38-4f2c766658b8', id FROM tags WHERE name = 'vermouth';

-- Recipe: Baked Lemon Cream Fish
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('819e527f-32db-49ca-b8a6-9c252d324791', 'Baked Lemon Cream Fish', '', 'food', 'western', '- 600g fish fillets
- 60g unsalted butter
- 1/2 cup cooking cream
- 3 cloves garlic, minced
- 30ml dijon mustard
- 1 lemon, juiced
- Salt & pepper
- fresh parsley, chopped (garnish)
- 1 lemon, cut into wedges (garnish)', '- Preheat oven to 200°C (all oven types).
- Place fish in a baking dish - ensure the fish isn''t crammed in too snugly. See video or photos in post. Sprinkle both sides of fish with salt and pepper.
- Place butter, cream, garlic, mustard, lemon juice, salt and pepper in a microwave proof jug or bowl. Microwave in 2 x 30 sec bursts, stirring in between, until melted and smooth.
- Sprinkle fish with shallots, then pour over sauce.
- Bake for 10 - 12 minutes, or until fish is just cooked. Remove from oven and transfer fish to serving plates. Spoon over sauce, and garnish with parsley and lemon wedges if using.', '### Notes
- Great served with mashed potatoes or rice to soak up the sauce
- Used frozen whiting fillets from Aldi for this and it worked great, note this goes directly against the advise on recipetineats website, so YMMV.', '**Name:** recipe tin eats
**URL:** https://www.recipetineats.com/baked-fish-with-lemon-cream-sauce/
**Type:** copy
**Modifications:** ratios, fish type', 'croach', '2023-11-22 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('fish');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '819e527f-32db-49ca-b8a6-9c252d324791', id FROM tags WHERE name = 'fish';
INSERT OR IGNORE INTO tags (name) VALUES ('lemon');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '819e527f-32db-49ca-b8a6-9c252d324791', id FROM tags WHERE name = 'lemon';

-- Recipe: Bastard Beans
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('8c28aaf3-fc36-4ad6-beee-6a50b74c72ce', 'Bastard Beans', '', 'food', 'mexican', '- 1 onion, grated (important)
- 10g garlic, fincely minced
- 2x400g beans
- olive oil
- water

### spices
- 1 tsp cumin ground
- 1 tsp corriander ground
- 1 tsp sweet paprika
- 1 tsp tsp oregano
- 1 tsp salt', '- Combine spices
- Heat oil in a pot
- Add onion and garlic paste, cook down
- Add the spices to the onion and garlic, waking up the spices. Add a splash of water if its looking dry like it might burn
- add the beans and cook for a minute
- cover the beans with water and cook for further 20 mins or so
- once the water has reduced and the beans are nice and soft, mash them to your preferred level of mashedness (i like it like 75% mashed)
- adjust thickness by cooking longer, or by adding more water', '### Notes
- mexican 3 bean mix is a favourite
- black beans i like as a mixer', '**Name:** string
**URL:** string
**Type:** string
**Modifications:** string', 'string', '2025-07-27 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('beans');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '8c28aaf3-fc36-4ad6-beee-6a50b74c72ce', id FROM tags WHERE name = 'beans';
INSERT OR IGNORE INTO tags (name) VALUES ('with rice');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '8c28aaf3-fc36-4ad6-beee-6a50b74c72ce', id FROM tags WHERE name = 'with rice';

-- Recipe: Beef And Broccoli Noodles
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('be51bcad-3ec6-4e05-8e28-4fa185801e22', 'Beef And Broccoli Noodles', '', 'food', 'chinese', '- 500g beef quick cooking beef, thinly sliced e.g. scotch, rump
- 1 1/2 tbsp peanut or vegetable oil
- 2 garlic cloves, finely chopped
- 1 onion, sliced
- 2 large heads of broccoli , broken into small florets
- 400g egg noodles (hokkien, lo mein, singapore - from the fridge)

### sauce
- 1/2 cup / 125 ml water
- 1 tbsp cornflour / cornstarch
- 2 tbsp dark soy sauce
- 2 tbsp light soy sauce (bumped from OG of 1.5, test going)
- 1 1/2 tbsp Chinese cooking wine/Shoasing wine
- 1 tsp white sugar (omit if using Mirin)
- 1/8 tsp Chinese five spice powder (not critical)
- 1/2 tsp sesame oil (not critical)
- 1/4 tsp pepper (white or black)

### optional garnishes
- sesame seeds
- chopped green onions', '### prep components
- Make the sauce by placing water and cornflour in a bowl and mixing. Then add remaining ingredients and mix.
- Place beef in a bowl and add 1 1/2 tbsp of the Sauce you just made. Mix it up.
- Bring a large pot of water to the boil. Add broccoli, cook for 1 minute. Add noodles then after 15 seconds, use a wooden spoon to separate the noodles then immediately drain (don''t have the noodles in the water for more than 1 minute)."

### do the stir fry
- Heat oil in a large skillet over high heat.
- Add garlic, quickly stir. Add onion and cook for 1 minute until it''s tinged with brown.
- Add beef and cook until it changes from red to brown.
- Add noodles, broccoli and sauce. Toss together for 1 1/2 - 2 minutes or until Sauce thickens and coats the noodles. Apply garnish.', '', '**Name:** recipe tin eats
**URL:** https://www.recipetineats.com/chinese-beef-broccoli-noodles/
**Type:** copy
**Modifications:** a few', 'croach', '2023-10-20 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('beef');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'be51bcad-3ec6-4e05-8e28-4fa185801e22', id FROM tags WHERE name = 'beef';
INSERT OR IGNORE INTO tags (name) VALUES ('noodles');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'be51bcad-3ec6-4e05-8e28-4fa185801e22', id FROM tags WHERE name = 'noodles';

-- Recipe: Beef Mince Bulgogi
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('4deada78-7ead-4fcd-b097-039466a8b967', 'Beef Mince Bulgogi', '', 'food', 'korean', '### main part of recipe
- 4 spring onions
- 1 large brown onion
- ~400g of mushrooms, whatever kind (shitake, button, oyster, dried, some mix of these)
- 1 green chili (optional, garnish)
- 1 red capsicum (optional), for colour
- 1 green capsicum (optional), for colour
- 1kg beef mince
- 4 tbsp neutral-tasting oil
- 1 tbsp toasted sesame oil
- Generous pinch of toasted sesame seeds, to garnish

### bulgogi marinade section
- 6 tbsp soy sauce
- 4 tbsp sugar
- 2 tbsp mirin
- 5 tbsp oyster sauce
- 4 tbsp minced garlic
- 2 tbsp toasted sesame oil
- Black pepper, to taste
- Small pinch of MSG

### TODO: need a way to add sections/headings to ingredients and instructions', '### PREP VEGETABLES
- Thinly slice the green onions, separating the whites and greens.
- Finely dice half an onion and mushrooms. Thinly slice the chilies (if using).

### MAKE BULGOGI MARINADE
- In a small container, mix together the soy sauce, sugar, mirin, oyster sauce, minced garlic, sesame oil, black pepper to taste, and a pinch of MSG (if using).
- To a large mixing bowl, add the ground beef, mushrooms, and marinade. Give it a good mix.

### COOK BULGOGI
- In a wok (or pan), heat the oil (2 tbsp) over medium-high heat. Once it gets nice and hot, add the white parts of the green onions and onion. Saute for 3 minutes or until they start to go brown.
- Add the beef mixture and cook for 7 to 8 minutes or until most of the liquid has evaporated. Make sure to break up the beef and keep stirring it.
- Turn the heat off. Add the chili peppers (if using), green onions, and sesame oil (1/2 tbsp). Give it a final mix.
- To serve, divide the rice into serving bowls and top with bulgogi. Garnish with some extra green onions and sesame seeds. Serve with salad greens or kimchi. Enjoy~!', '### Notes
- Serve with rice, add an egg or some greens.', '**Name:** Aaron & Claire
**URL:** https://aaronandclaire.com/ground-beef-bulgogi/
**Type:** copy
**Modifications:** minimal', 'croach', '2023-09-24 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('mince');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '4deada78-7ead-4fcd-b097-039466a8b967', id FROM tags WHERE name = 'mince';
INSERT OR IGNORE INTO tags (name) VALUES ('rice');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '4deada78-7ead-4fcd-b097-039466a8b967', id FROM tags WHERE name = 'rice';

-- Recipe: Bolognese Sauce
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('1cb3a15f-1f72-416f-b863-83912a762cdf', 'Bolognese Sauce', '', 'food', 'italian', '- 6 carrots (~640g) chopped
- 6 celery (~690g) chopped
- 3 medium onions (~280g) chopped
- 80ml EVOO
- 1kg beef mince (10% fat)
- 1kg pork mince (10% fat) (mix lean and regular)
- 35g garlic
- 2L passata
- 2 cans crushed tomatos
- 300g tomato paste
- 1.5tsp white sugar
- 1 good pinch MSG
- 400ml white wine
- dried italian herb mix
- salt and pepper', '- Place a large pot on medium heat, add EVOO and once warm, add onion. Cook until soft.
- Add carrots and celery. Leave to cook for 10 minutes.
- Add 200ml of the wine and leave to simmer on low until the alcohol evaporates.
- Add the mince and break it down before seasoning with salt and pepper.
- Leave mince to brown and excess water to evaporate.
- Mix in the remaining 200ml of wine.
- After allowing the wine to evaporate again, and garlic and cook for 30s.
- Mix through the passata, tomato paste and peeled tomatoes.
- Cook the sauce down to desired consistency (should take around 90mins).
- Season with salt, msg and dried herbs.
- You are done and may eat. Optionally, add some water and cook the sauce down again. Repeat for up to 4 or 5 hours of total cooking time.', '### Notes
- Makes about 5.5kg of sugo
- stir through a little bit of milk or cream before serving to make it delicious

### Next
- Less tomato paste', '**Name:** Vincenzo''s plate
**URL:** https://www.vincenzosplate.com/authentic-bolognese-sauce/
**Type:** string
**Modifications:** string', 'string', '2024-09-05 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('pasta');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '1cb3a15f-1f72-416f-b863-83912a762cdf', id FROM tags WHERE name = 'pasta';
INSERT OR IGNORE INTO tags (name) VALUES ('tomato');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '1cb3a15f-1f72-416f-b863-83912a762cdf', id FROM tags WHERE name = 'tomato';
INSERT OR IGNORE INTO tags (name) VALUES ('mince');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '1cb3a15f-1f72-416f-b863-83912a762cdf', id FROM tags WHERE name = 'mince';
INSERT OR IGNORE INTO tags (name) VALUES ('sauce');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '1cb3a15f-1f72-416f-b863-83912a762cdf', id FROM tags WHERE name = 'sauce';

-- Recipe: Chef John's Meatballs
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('0fed6405-4a03-4949-b0cc-a9c6cebec211', 'Chef John''s Meatballs', '', 'food', 'italian', '- ⅓ cup plain bread crumbs
- ½ cup milk
- 2 tablespoons olive oil
- 1 onion, diced
- 500g ground beef
- 500g ground pork
- 2 eggs
- 2 tablespoons grated Parmesan cheese
- ¼ bunch fresh parsley, chopped
- 3 cloves garlic, crushed
- 2 teaspoons salt
- 1 teaspoon ground black pepper
- 1 teaspoon dried Italian herb seasoning
- ½ teaspoon red pepper flakes (optional)', '- Cover a baking sheet with foil and spray lightly with cooking spray. Soak bread crumbs in milk in a small bowl for 20 minutes.
- Meanwhile, heat olive oil in a skillet over medium heat. Add onion; cook and stir until onion has softened and turned translucent, about 5 minutes. Reduce heat to low and continue cooking and stirring until onion is very tender, about 15 minutes more.
- Gently stir beef and pork together in a large bowl. Add onions, bread crumb mixture, eggs, Parmesan cheese, parsley, garlic, salt, black pepper, Italian herb seasoning, and red pepper flakes; mix together using a rubber spatula until combined. Cover and refrigerate for about one hour.
- Preheat the oven to 220 degrees C.
- Form mixture into balls about 1 1/2 inches in diameter; arrange in a single layer on the prepared baking sheet.
- Bake in the preheated oven until browned and cooked through, 15 to 20 minutes.', '### Notes
- Great to serve with a red sauce pasta, or as is.', '**Name:** Chef John
**URL:** https://www.allrecipes.com/recipe/220854/chef-johns-italian-meatballs/
**Type:** copy
**Modifications:** minimal', 'croach', '2023-09-14 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('mince');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '0fed6405-4a03-4949-b0cc-a9c6cebec211', id FROM tags WHERE name = 'mince';

-- Recipe: Chilli Wine Garlic Prawn Pasta
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('b601d7e0-be25-452e-9d70-5d05034f5abb', 'Chilli Wine Garlic Prawn Pasta', '', 'food', 'fusion', '- 1kg prawns
- 360g linguine or similar
- 3 shallot, diced or half a white onion
- 2 red chilli, deseeded and diced
- 400 grams cherry tomatoes, halved
- 5 cloves garlic, finely sliced
- 1 cup white wine,
- 2 tbsp unsalted butter
- zest of 2 lemon
- 2 tbsp fresh chopped parsley
- 40 grams grated parmesan cheese plus extra for serving if desired', '- Bring a large pot of salted water to boil for cooking the pasta.
- Add the butter to a frying pan over medium heat and, if using raw prawns, add them to the pan and cook until they just turn pink. Remove the prawns and set aside, leaving the butter in the pan.
- Add the shallot and chilli to the pan with the butter and saute. Meanwhile, put the pasta on to cook.
- When the shallot begins to soften (2-3 minutes) add the garlic and tomatoes to the pan with a pinch of salt and pepper. Stir together then pour over the wine.
- Allow the sauce to simmer for about 5 minutes until the pasta is cooked. It should be slightly al dente (use the lower cooking time recommended on the packet).
- Add the prawns into the pan with the chilli and tomato, give everything a stir and add the pasta. (I prefer to use tongs to transfer the pasta directly into the sauce. This also adds some of the pasta water which will emulsify the sauce and help it stick to the pasta. If you prefer to drain the pasta first in a colander, make sure to keep a tablespoon of the cooking water and add it to the sauce with the pasta.)
- Sprinkle over the parsley, parmesan and lemon zest and toss everything together with tongs before serving.
- Serve with an extra sprinkle of parmesan and a side of salad.', '### Notes
- for wine i''ve used a drinking quality NZ sav blanc and it was excellent
- if using frozen prawns, then instead start the cooking by just cooking the prawns on their own to take the chill off them and remove lots of ice. Once they are no longer icy (but not fully cooked yet) remove the prawns, and put them aside to incorporate as guided in instructions.
- also works by using frozen fish such as barramundi, basa etc
- also works with the addition of something green, like broccoli or asparagus. For brocolli, par boil first, then add in when u add in the fish

### Next
- pasta quantity is a bit of a guess, measure and see what works', '**Name:** Carrie''s Kitchen
**URL:** https://carriecarvalho.com/chilli-prawn-linguine/', '', '2023-12-19 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('pasta');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'b601d7e0-be25-452e-9d70-5d05034f5abb', id FROM tags WHERE name = 'pasta';
INSERT OR IGNORE INTO tags (name) VALUES ('prawn');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'b601d7e0-be25-452e-9d70-5d05034f5abb', id FROM tags WHERE name = 'prawn';

-- Recipe: Creamy Chicken Mushroom Pasta
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('65dd1b8f-c4f3-4c55-bd07-9aed6da3cab9', 'Creamy Chicken Mushroom Pasta', '', 'food', 'western', '- 1kg chicken breast, cubed
- 500g mushrooms, sliced
- 280g fresh spinach leaves
- 1 340ml carnation creamy evaporated milk
- 500g pasta, farfalline, farfalle, penne etc
- 3tbsp olive oil
- 2 onions
- 5 garlic cloves
- 1 teaspoon dried thyme
- 2 teaspoons paprika
- 1L chicken broth
- salt
- pepper
- 50g grated parmesan', '- Heat 2tbsp olive oil in a large pot on medium heat. Add chicken, salt and pepper. Cook through and set aside.
- Add 1tbsp olive oil and the onion into the pot and stir. Cook down for 1-2 minutes.
- Add mushrooms and garlic, and stir to incorporate with the onion. Season with salt and pepper to taste as well as thyme and paprika. Stir to evenly season.
- Add chicken broth and evaporated milk to the pot and stir. Bring to a boil, then add the pasta.
- Cook according to package instructions, being sure to stir every 1-2 minutes to keep the pasta from clumping together. (Cook time may be a little longer in this recipe than when the pasta is boiled in water.)
- When the farfalle pasta is al dente, add the spinach and chicken and stir until the spinach cooks down and incorporates.
- Add parmesan and stir until it''s well-incorporated and you''re left with a smooth sauce.
- Top off with extra parmesan and serve.', '', '**Name:** tasty
**URL:** https://tasty.co/recipe/one-pot-chicken-and-mushroom-pasta
**Type:** string
**Modifications:** string', 'meggles', '2025-04-25 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('chicken');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '65dd1b8f-c4f3-4c55-bd07-9aed6da3cab9', id FROM tags WHERE name = 'chicken';
INSERT OR IGNORE INTO tags (name) VALUES ('mushroom');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '65dd1b8f-c4f3-4c55-bd07-9aed6da3cab9', id FROM tags WHERE name = 'mushroom';
INSERT OR IGNORE INTO tags (name) VALUES ('pasta');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '65dd1b8f-c4f3-4c55-bd07-9aed6da3cab9', id FROM tags WHERE name = 'pasta';
INSERT OR IGNORE INTO tags (name) VALUES ('one pot');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '65dd1b8f-c4f3-4c55-bd07-9aed6da3cab9', id FROM tags WHERE name = 'one pot';
INSERT OR IGNORE INTO tags (name) VALUES ('dutch oven');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '65dd1b8f-c4f3-4c55-bd07-9aed6da3cab9', id FROM tags WHERE name = 'dutch oven';

-- Recipe: Easy Chicken Curry
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('284014c7-73fd-46b7-bbed-520ea6f37ed5', 'Easy Chicken Curry', '', 'food', 'indian', '- 1kg chicken breast/thighs, cut into cubes
- 2 onions
- 3 garlic cloves, grated
- 1 knob ginger, grated
- 3 tablespoons vegetable oil
- 2 tablespoon curry powder
- 1.5 teaspoons salt
- 1/2 teaspoon black pepper
- 4 large tomatoes, chopped
- 1 can (400g) coconut cream/milk
- 180ml water
- ½ teaspoon garam masala
- Coriander (chopped for garnish)', '- In a large pan heat oil. Add chopped onions and sauté until golden, over medium-low heat. stir occasionally.
- When onions are golden add garlic and ginger. Stir and cook for 2 minutes. Add chopped tomatoes and cook, stirring occasionally, until soft. Add curry powder and cook, stirring, for 2 minutes. Add chicken breast, season with salt and pepper, cook for 5 minutes, until golden. Cover the pan and cook for 2-3 minutes more. Add water, cover and cook 10 minutes.
- Add coconut milk/cream, season with garam masala and salt, stir, and cook until lightly thickens.
- Serve with naan bread or rice.', '', '**Name:** The Cooking Foodie
**URL:** https://www.thecookingfoodie.com/recipe/Quick-and-Easy-Chicken-Curry-Recipe
**Type:** copy
**Modifications:** minimal', 'snapper', '2023-09-24 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('curry');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '284014c7-73fd-46b7-bbed-520ea6f37ed5', id FROM tags WHERE name = 'curry';
INSERT OR IGNORE INTO tags (name) VALUES ('chicken');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '284014c7-73fd-46b7-bbed-520ea6f37ed5', id FROM tags WHERE name = 'chicken';

-- Recipe: Ez Green Chicken Curry
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('179fc586-7018-491f-9755-4b6618ffd0c8', 'Ez Green Chicken Curry', '', 'food', 'thai', '### vegetables
- ~600g assorted chopped vegetables (think snow peas, eggplant, asparagus, carrot, green beans, zuchinni) etc
- 150g onion
- 20g garlic
- 10g ginger
- 1.2kg chicken breast, diced (or anything else)
- 500ml chicken stock
- 400ml can coconut milk (full fat, ayam is good)
- 195g ayam thai green curry paste jar
- 2tsp sugar
- 6tsp fish sauce', '- In oil, cook onions until soft/transulscent.
- Add curry paste, ginger, garlic, let simmer and cook for a bit
- Add the chicken and cook it in the curry paste mixture, tossing a bit. Just a few minutes to give some colour and sear the outside.
- Add in stock, coconut milk, sugar and fish sauce.
- Now we will simmer for 20 minutes and we will be ready to eat. Simmer the vegetables for some of that 20minutes depending on how thick you cut them and what they are. Or just dump them all in and forget - it will be fine.', '### Next
- add lime, lemongrass paste or frozen lemongrass chopped up, and fakkir lime leaves?', '**Name:** recipetineats
**URL:** https://www.recipetineats.com/thai-green-curry/#h-the-best-green-curry-paste
**Type:** made simpler, changed ratios, specified brands, better macros
**Modifications:** string', 'string', '2025-03-28 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('curry');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '179fc586-7018-491f-9755-4b6618ffd0c8', id FROM tags WHERE name = 'curry';
INSERT OR IGNORE INTO tags (name) VALUES ('chicken');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '179fc586-7018-491f-9755-4b6618ffd0c8', id FROM tags WHERE name = 'chicken';
INSERT OR IGNORE INTO tags (name) VALUES ('jar');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '179fc586-7018-491f-9755-4b6618ffd0c8', id FROM tags WHERE name = 'jar';
INSERT OR IGNORE INTO tags (name) VALUES ('with rice');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '179fc586-7018-491f-9755-4b6618ffd0c8', id FROM tags WHERE name = 'with rice';

-- Recipe: Ez Guilin Chicken
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('1283c346-4d31-480d-ac67-bf8980d5d7e4', 'Ez Guilin Chicken', '', 'food', 'chinese', '- 1 large onion
- 4 cloves of garlic, minced or paste
- knob of ginger, minced or paste
- 1.2kg chicken thighs cut into bite sized pieces
- 2 tbsp oyster sauce
- 2 tbsp guilin sauce (lee kum kee brand)', '- cook down onions until translucent
- add in ginger and garlic and cook for a minute
- add in chicken, stiring
- once the chicken is half cooked, add in the sauces
- cook until chicken is done and sauce is appropriately thick. Add a splash water if you want to loosen it up
- nice as meal prep served with rice and some veg', '', '**Name:** OG
**URL:** string
**Type:** string
**Modifications:** string', 'string', '2025-09-10 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('chicken');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '1283c346-4d31-480d-ac67-bf8980d5d7e4', id FROM tags WHERE name = 'chicken';
INSERT OR IGNORE INTO tags (name) VALUES ('with rice');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '1283c346-4d31-480d-ac67-bf8980d5d7e4', id FROM tags WHERE name = 'with rice';
INSERT OR IGNORE INTO tags (name) VALUES ('easy');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '1283c346-4d31-480d-ac67-bf8980d5d7e4', id FROM tags WHERE name = 'easy';
INSERT OR IGNORE INTO tags (name) VALUES ('spicy');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '1283c346-4d31-480d-ac67-bf8980d5d7e4', id FROM tags WHERE name = 'spicy';

-- Recipe: Ez Mapo Tofu
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('208dc31b-45eb-4ec8-89ad-47de72408983', 'Ez Mapo Tofu', '', 'food', 'chinese', '- neutral oil
- 1 brown onion, sliced finely
- 2 cloves garlic, grated/sliced/chopped
- 1 knob of ginger, grated/sliced/chopped
- 500g pork mince
- 900g tofu, cut into cubes
- spring onions (garnish, optional)

### sauces
- 113g of Lee Kum Kee mapo tofu sauce (usually do half the 226g jar, so 113g)
- 1 tablespoon oyster sauce
- 2 teaspoons soy sauce', '- Combine the sauces
- Get some oil hot, saute the onions down a bit
- Add the ginger and garlic and cook for 30s
- Add the mince in and break it up and cook it almost all through
- Add in the sauce and finish cooking the mince. Add water if things are getting too dry to loosen it up.
- Gently add in the tofu and little it warm through a bit. Try not to break it much before it is enjoyed.', '### Notes
- Serve with rice

### Next
- the added oyster sauce and soy sauce was a guess at what i did, validate this and iterate', '**Type:** pimping/original/research inspired', 'croach', '2023-11-25 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('tofu');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '208dc31b-45eb-4ec8-89ad-47de72408983', id FROM tags WHERE name = 'tofu';
INSERT OR IGNORE INTO tags (name) VALUES ('pork');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '208dc31b-45eb-4ec8-89ad-47de72408983', id FROM tags WHERE name = 'pork';
INSERT OR IGNORE INTO tags (name) VALUES ('spicy');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '208dc31b-45eb-4ec8-89ad-47de72408983', id FROM tags WHERE name = 'spicy';
INSERT OR IGNORE INTO tags (name) VALUES ('rice');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '208dc31b-45eb-4ec8-89ad-47de72408983', id FROM tags WHERE name = 'rice';
INSERT OR IGNORE INTO tags (name) VALUES ('jar');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '208dc31b-45eb-4ec8-89ad-47de72408983', id FROM tags WHERE name = 'jar';
INSERT OR IGNORE INTO tags (name) VALUES ('with rice');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '208dc31b-45eb-4ec8-89ad-47de72408983', id FROM tags WHERE name = 'with rice';

-- Recipe: Firecracker Chicken
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('b0d94d13-6f4e-487a-b98c-2effc5efe8a8', 'Firecracker Chicken', '', 'food', 'fusion', '- 500g chicken mince
- 10h chilli oil
- 30g honey
- 60g frank red hot buffalo sauce
- 30g rice vingegar
- 1tsp garlic powder
- 1 onion, chopped small
- 3 birds eye chillis
- 2 garlic gloves
- equivalent ginger', '- create the sauce by mixing honey, franks hot sauce and vinegar
- Heat a pan with some oil
- cook the onions for a minute or two, until a bit soft
- add the chilli oil and cook the garlic, ginger and chilllis for 20 seconds, just take the edge off
- add the chicken and cook
- once chicken is done, add the sauce and reduce it down to desired thickness', '### Notes
- serve with rice, broccoli and a boiled egg', '**Name:** KindaHealthyRecipes
**URL:** https://masonfit.com/low-carb-firecracker-ground-chicken/
**Type:** string
**Modifications:** string', 'string', '2024-06-08 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('mealprep');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'b0d94d13-6f4e-487a-b98c-2effc5efe8a8', id FROM tags WHERE name = 'mealprep';
INSERT OR IGNORE INTO tags (name) VALUES ('spicy');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'b0d94d13-6f4e-487a-b98c-2effc5efe8a8', id FROM tags WHERE name = 'spicy';
INSERT OR IGNORE INTO tags (name) VALUES ('mince');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'b0d94d13-6f4e-487a-b98c-2effc5efe8a8', id FROM tags WHERE name = 'mince';
INSERT OR IGNORE INTO tags (name) VALUES ('with rice');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'b0d94d13-6f4e-487a-b98c-2effc5efe8a8', id FROM tags WHERE name = 'with rice';

-- Recipe: Gyudon
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('e6bc7920-0c01-48a8-a49c-77b22cbeffa3', 'Gyudon', '', 'food', 'japanese', '- 750g thinly sliced beef strips
- 2 onions, thinly sliced into long strips
- 3 cups boiling water
- 2 teaspoons dashi powder
- 90ml soy sauce
- 90ml mirin
- 45ml sake
- 2 tablespoons sugar
- 5 eggs, lightly beaten
- pinch MSG
- rice
- spring onion, sliced (garnish, optional)
- pickled ginger (garnish, optional)', '- Put the rice cooker on in advance.
- In a pot combine water, dashi, soy sauce, mirin, sake, msg and sugar. Stir until disolved.
- Once simmering, add onions and cook until soft.
- Next, add the meat, stir occasionally and cook. It''ll be done quick, couple minutes.
- Next, slowly pour the eggs over the top, distributing evenly. Chuck the lid on so the eggs cook and steam for a minute.
- Serve over rice and eat immediately. Garnish with spring onion and/or pickled ginger.', '### Notes
- You might choose to replace pouring the egg in, with instead serving the meal with a fried egg on top for each bowl.
- I think this might be great with lamb slices too
- For the sliced meat, easiest way is to buy the hot pot/sukiyaki/bulgogi meat from freezer section of asian grocer, but has higher fat content. Healthier to buy leaner meat and slice yourself.', '**Type:** research inspired/amalgamation', 'croach', '2023-11-29 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('beef');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'e6bc7920-0c01-48a8-a49c-77b22cbeffa3', id FROM tags WHERE name = 'beef';
INSERT OR IGNORE INTO tags (name) VALUES ('rice');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'e6bc7920-0c01-48a8-a49c-77b22cbeffa3', id FROM tags WHERE name = 'rice';
INSERT OR IGNORE INTO tags (name) VALUES ('with rice');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'e6bc7920-0c01-48a8-a49c-77b22cbeffa3', id FROM tags WHERE name = 'with rice';

-- Recipe: Healthy Curry
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('7c7530c3-9bcb-4db6-97b6-c7630c5ab5fa', 'Healthy Curry', '', 'food', 'indian', '### Protein & marinade
- 10g tbsp oil
- 1.2kg chicken thighs, each cut into 6
- 1 tsp tumeric powder
- 1 tsp kashmiri chilli powder
- 0.5 tsp cumin powder
- 0.5 tsp corriander powder
- half a lemon

### Base
- 2 tbsp oil
- 400g onion, finely sliced
- 100g carrot, chopped small
- 150g green capsicum, roughly chopped
- 30g ginger, minced or grated
- 30g garlic, minced or grated

### Spices for oil
- 2 tbsp oil
- 1 bay leaf
- 5 cloves
- 5 green cardamon
- 0.25 tsp mustard seeds
- 0.25 tsp cumin seeds

### Powder additions
- 1.5 tsp tumeric
- 1 tsp kashmiri chilli powder
- 3 tsp corriander powder
- 2 tsp cumin powder
- 1 tsp garam masala

### Veg
- 240g (1 drained can) chickpeas
- 260g green beans
- NEED MORE', '- Combine the chicken thighs with the marinade spices and let marinade, while you do the next things, or overnight.
- Start the base gravy, by first adding 2tbsp oil and onions to a large pot. Let cook down until nice and soft, atleast 15min.
- Add ginger and garlic and cook until fragrant, just a minute or so.
- Add the carrot and capsicum and cook down a bit, then add the tomato and cook it all down. Add some water if any risk of burning at any stage.
- Once it''s cooked down, transfer to a blender and blend. Set aside for later.
- Now back in the same (now empty) pot, add another 2 tbsp of oil. Once hot, add the whole spices and cook for a minute.
- Add the chicken and brown it all around in the fragrant oil.
- Once the chicken is sealed on the outside, add in the gravy sauce and the powder additions.
- Add the veg and cook it and its done.', '### Notes
- string

### Next
- more Veg
- more tomato
- tomato paste?
- put the powder additions earlier? so they get fried off too?', '**Name:** string
**URL:** string
**Type:** string
**Modifications:** string', 'string', '2025-10-22 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('curry');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '7c7530c3-9bcb-4db6-97b6-c7630c5ab5fa', id FROM tags WHERE name = 'curry';
INSERT OR IGNORE INTO tags (name) VALUES ('with recipe');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '7c7530c3-9bcb-4db6-97b6-c7630c5ab5fa', id FROM tags WHERE name = 'with recipe';

-- Recipe: Honey Pepper Stirfry
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('dd8fdffc-12c7-4736-a0d0-bfd6d4b7e199', 'Honey Pepper Stirfry', '', 'food', 'chinese', '### sauce
- 100ml soy sauce
- 60g honey
- 70g tbsp Oyster sauce (sub Hoisin)
- 3 tbsp Chinese cooking wine , Mirin or dry sherry (Note 2)
- 2 tbsp water
- 1 tsp coarsely crushed black pepper (or 1/2 tsp normal ground black pepper)

### Stir Fry
- 2 garlic cloves, minced
- 1 onion, sliced
- 500g thinly sliced beef steak
- 1 carrot, julienned
- 1 head broccoli, cut to small floretts
- 450g tofu, cut into cubes', '- Mix the Sauce ingredients in a bowl.
- Steam the broccoli and carrot for a couple mins in the microwave
- Heat some oil in a wok or large heavy based skillet over high heat until it is smoking.
- Add the onion and garlic and cook for 1 minute or until the onion becomes translucent. Keep it moving so the garlic doesn''t burn.
- Add the beef and stir fry for 1 minute until just cooked to your liking, then remove into bowl.
- Turn the heat down to medium high, pour in Sauce - it will start simmering very quickly! Let it cook for 1 minute or so until it becomes syrupy - the bubbles will be larger and caramel colour.
- Add the tofu in and let it warm
- Add the vegetables, beef and onion back into the wok, along with any juices pooled on the plate. Toss in the sauce until just warmed through - 1 minute at most. Don''t overcook the beef - that would be tragic!
- Serve with rice', '### Notes
- would work with pork or chicken too', '**Name:** recipetineats
**URL:** string
**Type:** string
**Modifications:** string', 'string', '2024-04-06 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('beef');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'dd8fdffc-12c7-4736-a0d0-bfd6d4b7e199', id FROM tags WHERE name = 'beef';
INSERT OR IGNORE INTO tags (name) VALUES ('stirfry');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'dd8fdffc-12c7-4736-a0d0-bfd6d4b7e199', id FROM tags WHERE name = 'stirfry';
INSERT OR IGNORE INTO tags (name) VALUES ('honey');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'dd8fdffc-12c7-4736-a0d0-bfd6d4b7e199', id FROM tags WHERE name = 'honey';
INSERT OR IGNORE INTO tags (name) VALUES ('pepper');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'dd8fdffc-12c7-4736-a0d0-bfd6d4b7e199', id FROM tags WHERE name = 'pepper';

-- Recipe: Japanese Chicken Curry
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('e986a7a6-104c-446c-8e0b-2de6bc496cf0', 'Japanese Chicken Curry', '', 'food', 'string', '- 670g onions, sliced
- 2 carrots (6.7 oz, 190 g)
- 3 Yukon gold potatoes (15 oz, 432 g)
- 1 tsp ginger (grated)
- 2 cloves garlic
- ½ apple (6 oz, 170 g)
- 1½ lb boneless, skinless chicken thighs (see Notes for substitutions)
- freshly ground black pepper

### curry sauce
- 1½ Tbsp neutral oil (for cooking)
- 4 cups low sodium chicken stock/broth (the curry already salty, can always add more salt later to taste)
- 1 Tbsp honey
- 1 Tbsp soy sauce
- 1 Tbsp ketchup
- 200g Japanese curry roux (1 packet)', '- Heat 1½ Tbsp neutral oil in a large pot over medium heat and add the onion.
- Sauté the onions, stirring occasionally, until they become translucent and tender, about 5 minutes. Don''t stir them too often because they won''t develop a golden color. Cooked onions add amazing flavor, so don''t skip this step. If you have extra time, definitely sauté the onions an additional 5 minutes to add more color and flavor.
- Add the minced garlic (I pass it through a garlic press for a finer texture) and grated ginger and mix well together.
- Add the chicken and cook, stirring frequently, until it''s no longer pink on the outside. If the onions are getting too brown, reduce the heat to medium low temporarily.
- Add 4 cups chicken stock/broth.
- Add the grated apple, 1 Tbsp honey, 1 Tbsp soy sauce, and 1 Tbsp ketchup (or any condiment you choose to add). Please read my blog post for details.
- Add the carrots and Yukon gold potatoes (if you''re using russet potatoes, add them later in the last 15-20 minutes of cooking). The broth should barely cover the ingredients. Don''t worry; we don''t want too much liquid here, and additional liquid will be released from the meat and vegetables.
- Simmer, covered*, on medium-low heat for 15 minutes, stirring occasionally. *Simmer uncovered if the ingredients are completely submerged in the broth.
- Once boiling, use a fine-mesh strainer to skim the scum and foam from the surface of the broth.
- Continue to cook covered until a wooden skewer goes through the carrots and potatoes.
- To Add the Curry Roux
- Turn off the heat. From 1 package Japanese curry roux, put 1-2 cubes in a ladleful of cooking liquid. Slowly let it dissolve with a spoon or chopsticks and stir into the pot to incorporate. Repeat with the rest of the blocks, 2 cubes at a time. Tip - I use 1 full-sized box of store-bought curry roux, which is typically for 8-12 servings (be careful, as some brands offer a smaller box, which is 4 servings). With my homemade curry roux, I typically use 6-7 cubes for 8 servings (about 80% of the curry roux mixture if it hasn''t solidified yet).
- Simmer, uncovered, on medium-low heat, stirring frequently, until the curry becomes thick, about 5-10 minutes. If your curry is too thick, you can add water to thin the sauce. When you stir, make sure that no roux or food is stuck to the bottom of the pot; otherwise, it may burn.', '### Notes
- string

### Next
- string', '**Name:** string
**URL:** string
**Type:** string
**Modifications:** string', 'string', '2023-11-27 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('with rice');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'e986a7a6-104c-446c-8e0b-2de6bc496cf0', id FROM tags WHERE name = 'with rice';
INSERT OR IGNORE INTO tags (name) VALUES ('curry');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'e986a7a6-104c-446c-8e0b-2de6bc496cf0', id FROM tags WHERE name = 'curry';

-- Recipe: Kimchi Chikki Stew
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('fd0365e8-ac0f-4e5c-bf53-991742d9d81f', 'Kimchi Chikki Stew', '', 'food', 'korean', '### initial addition
- 1 container kimchi (~700g), and don''t skimp on quality
- 1.5kg chicken thighs
- 3 green onions, sliced
- 1 medium onion, sliced
- 3 garlic cloves, sliced
- a packet enoki mushrooms
- 1 teaspoon kosher salt
- 2 teaspoons sugar
- 2 teaspoons gochugaru (Korean hot pepper flakes)
- 1 tablespoon gochujang (hot pepper paste)
- 1 teaspoon toasted sesame oil
- 3 cups of dashi stock (3 tsp in 750ml)

### others
- 1 tube (350g) of sundubu soft tofu
- 1 pack of firmer stew tofu
- optional - add more kinds of veg like asian greens, potatos, lots works here', '- Open the kimchi container and use some scissors to chop it up (or a knife/chopping board)
- Into a large pot (green ceramic), place all the initial additions
- Cover and cook for 15 minutes.
- Remove lid and add the tofu, then cooking for another 15 minutes
- serve with rice', '', '**Name:** string
**URL:** https://www.maangchi.com/recipe/kimchi-jjigae
**Type:** string
**Modifications:** based on but many changes', 'string', '2025-09-10 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('with rice');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'fd0365e8-ac0f-4e5c-bf53-991742d9d81f', id FROM tags WHERE name = 'with rice';
INSERT OR IGNORE INTO tags (name) VALUES ('spicy');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'fd0365e8-ac0f-4e5c-bf53-991742d9d81f', id FROM tags WHERE name = 'spicy';
INSERT OR IGNORE INTO tags (name) VALUES ('soup');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'fd0365e8-ac0f-4e5c-bf53-991742d9d81f', id FROM tags WHERE name = 'soup';
INSERT OR IGNORE INTO tags (name) VALUES ('stew');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'fd0365e8-ac0f-4e5c-bf53-991742d9d81f', id FROM tags WHERE name = 'stew';

-- Recipe: Lamb Burgers
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('e1b63d93-9d3b-4bdb-b3f5-3d0ef6b47064', 'Lamb Burgers', '', 'food', 'continental', '- 500g lamb mince
- 2 cloves garlic, minced
- 1 large egg
- 1/2 teaspoon cumin
- 1/2 teaspoon salt
- 1/4 teaspoon pepper
- zest of 1/2 lemon
- 1 teaspoon Rosemary
- 1 small onion, grated or very finely chopped', '- In a large mixing bowl, combine all the ingredients and mix until thoroughly combined.
- Using your hands, shape the mixture into four patties.
- The lamb patties should take around 4minutes each side, cooking closer to high so we get some nice browning.', '### Notes
- Burger inclusion inspo - cheese, red onion, tomato, lettuce, sauce (tomato, hot english mustard, other mustards, kewpie mayo all work)
- Also works as a greek flavours burger with tzatziki, feta etc...', '**Name:** The Big Man''s World
**URL:** https://thebigmansworld.com/lamb-burgers/
**Type:** copy', 'croach', '2023-12-26 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('lamb');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'e1b63d93-9d3b-4bdb-b3f5-3d0ef6b47064', id FROM tags WHERE name = 'lamb';
INSERT OR IGNORE INTO tags (name) VALUES ('burger');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'e1b63d93-9d3b-4bdb-b3f5-3d0ef6b47064', id FROM tags WHERE name = 'burger';

-- Recipe: Lamb Shanks Massaman Curry
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('2c4f2bad-0cec-458b-9c3a-b831a60efc3f', 'Lamb Shanks Massaman Curry', '', 'food', 'thai', '- 1.5kg lamb shanks (5 small, 4 medium, 2 - 3 large)
- 114g/ 4oz Maesri Massaman curry paste (1 can), or other brand (Note 1)
- 400ml coconut milk, full fat (Ayam brand is best, Note 3)
- 2 cups chicken stock/broth
- 1 onion
- 400g small potatoes
- 1 star anise
- 1 cinnamon stick
- Red chilli, finely sliced (garnish)
- Coriander (garnish)', '- Preheat oven to 180°C (160°C fan).
- Mix curry paste, coconut milk and stock in a baking dish. Add onion, potato, star anise, cinnamon and lamb.
- Turn shanks to coat in sauce, then cover with foil.
- Bake in oven for 2 hours. Remove foil, bake for a further 1 hour (small shanks) or 1.5 hrs (medium to large shanks), turning lamb twice to brown evenly, until meat is so tender it can easily be teased apart with 2 forks.
- Remove lamb onto plate. Carefully skim off excess fat off the surface (tilt dish, it''s easier) - I get about 1/3 cup. Mix sauce in baking dish - it should be reduced down to a syrupy thickness (Note 6).
- Serve lamb with sauce over jasmine rice, garnished with chilli and coriander.', '### Notes
- 1. This recipe is not suited to slow cooker, pressure cooker or stove. Achieving the intended flavour with such little effort requires the oven to caramelise the surface of the lamb and sauce (as well as reducing the sauce).
- 2. Massaman curry paste - best is Maesri brand, sold at most Woolworths & Coles in Australia. Ayam i tried and good (i put the whole can).
- 3. Coconut milk - the quality/flavour comes down to the % of the liquid that is actually coconut milk. Ayam is the highest at 89%, cheap ones can be as low as 45%.', '**Name:** recipe tin eats
**URL:** https://www.recipetineats.com/lamb-shanks-in-massaman-curry/
**Type:** copy
**Modifications:** minimal', 'croach', '2023-09-18 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('curry');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '2c4f2bad-0cec-458b-9c3a-b831a60efc3f', id FROM tags WHERE name = 'curry';
INSERT OR IGNORE INTO tags (name) VALUES ('lamb shanks');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '2c4f2bad-0cec-458b-9c3a-b831a60efc3f', id FROM tags WHERE name = 'lamb shanks';

-- Recipe: Mexican Meat Mix
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('65752f99-f4cc-4dd7-89ec-353afc973dc2', 'Mexican Meat Mix', '', 'food', 'mexican', '- 1 large onion, diced
- 3 cloves garlic, finely diced
- 500g extra lean beef mince
- 1 can mexican bean mix, drained and lightly rinsed
- 1 green capsicum, cut into medium sized peices
- 1 tomato
- 1 carrot, juillianned

### spice mix
- 1 tbsp sweet paprika
- 2 tsp cumin
- 2 tsp corriander
- 1 tsp salt
- 1 tsp extra hot chilli powder

### to season
- salt
- pepper', '- Combine spice mix
- Heat pan with medium heat
- In some oil add onion and cook for a minute or two. Then add garlic, cook for 30s.
- Add the mince meat, breaking it up and making sure its getting browned. Give it a little salt and pepper.
- Once most of the pink is gone, add in the beans. Combine and continue cooking till all the pink is gone.
- Add In the capsicum, carrot, and then tomato in installments, allowing them to soften a bit.
- Sprinkle over all the spice mix and combine well, adding water to get it slightly more wet than you want and letting it evaporate.', '### Notes
- serve with rice. garnish with conriander or parsley or soemthing?
- can add any kinds of fresh chillis when u add the garlic, or when u are adding veg towards the end

### Next
- wrote this recipe a bit baked so check the accuracy', '**Name:** original by davo', 'croach', '2024-08-20 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('beef');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '65752f99-f4cc-4dd7-89ec-353afc973dc2', id FROM tags WHERE name = 'beef';
INSERT OR IGNORE INTO tags (name) VALUES ('mince');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '65752f99-f4cc-4dd7-89ec-353afc973dc2', id FROM tags WHERE name = 'mince';
INSERT OR IGNORE INTO tags (name) VALUES ('beans');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '65752f99-f4cc-4dd7-89ec-353afc973dc2', id FROM tags WHERE name = 'beans';
INSERT OR IGNORE INTO tags (name) VALUES ('meal prep');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '65752f99-f4cc-4dd7-89ec-353afc973dc2', id FROM tags WHERE name = 'meal prep';
INSERT OR IGNORE INTO tags (name) VALUES ('lunch');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '65752f99-f4cc-4dd7-89ec-353afc973dc2', id FROM tags WHERE name = 'lunch';

-- Recipe: Minestrone
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('dee42d11-559c-4172-9482-90656b71fdfd', 'Minestrone', '', 'food', 'italian', '- 6 cups (1.5L) chicken stock
- 250g celery, diced
- 250g carrot, diced
- 250g onion, diced
- 4 cloves garlic, minced
- 500g turkey/chicken minced
- 500g pork mince
- 1 can (400g) diced tomatoes
- 1 can (400g) chickpeas, drained
- 1 can beans (e.g. red kidney, most will work), drained
- 1 cup passata
- 250g pasta

### spices
- 1 tbsp dried parsley
- 1.5 tsp dried basil
- 1 tsp dried oregano
- 0.5 tsp dried thyme
- 0.5 tsp salt
- 0.5 tsp black pepper

### optional garnish
- parmesan (optional garnish)
- lemon (optional garnish)
- chilli flakes (optional garnish)', '- Heat olive oil in a large saucepan or stock pot over medium heat.
- Add carrots, onion, and celery; cook for 2-3 minutes until onion is translucent.
- Add garlic and mince and cook until no longer pink.
- Add broth, diced tomatoes, tomato sauce, kidney beans, chickpeas, and spices.
- Reduce heat to low, and simmer for 20 minutes.
- Add pasta and cook for another 7-8 minutes until al dente.
- Garnish with parmesan, a squeeze of lemon, or red pepper flakes if desired.', '', '**Name:** string
**URL:** https://nutritionistmom.com/blogs/blog/high-protein-minestrone-soup?srsltid=AfmBOooHIvu5diRZIYGwelSVvXfbRB5I7GiILuxsf1XyVWy-PzgQWr-R
**Type:** string
**Modifications:** string', 'string', '2025-08-14 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('soup');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'dee42d11-559c-4172-9482-90656b71fdfd', id FROM tags WHERE name = 'soup';
INSERT OR IGNORE INTO tags (name) VALUES ('healthy');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'dee42d11-559c-4172-9482-90656b71fdfd', id FROM tags WHERE name = 'healthy';
INSERT OR IGNORE INTO tags (name) VALUES ('fibre');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'dee42d11-559c-4172-9482-90656b71fdfd', id FROM tags WHERE name = 'fibre';

-- Recipe: Oyakodon
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('779b1147-1777-4b0a-b188-c60af0b2a2a7', 'Oyakodon', '', 'food', 'japanese', '- 2 cups (500ml) dashi stock
- 60ml sake
- 30ml soy sauce, plus more to taste
- 30g sugar, plus more to taste
- 2 large onions (~300g), thinly sliced
- 600g boneless, skinless chicken thighs, chopped to bitesize pieces
- 4 green onions, sliced, divided in 2 (optional; see note)
- ~6 large eggs
- togarashi (japanese chilli powder) (optional; garnish/seasoning)
- rice', '### TODO: break up to more steps that are shorter?
- Combine dashi, sake, soy sauce, and sugar in a 10-inch skillet and bring to a simmer over high heat. Adjust heat to maintain a strong simmer. Stir in onion and cook, stirring occasionally, until onion is half tender, about 5 minutes. Add chicken pieces and cook, stirring and turning chicken occasionally, until chicken is cooked through and broth has reduced by about half, 5 to 7 minutes for chicken thighs or 3 to 4 minutes for chicken breast. Stir in half of scallions and all of mitsuba (if using), then season broth to taste with more soy sauce or sugar as desired. The sauce should have a balanced sweet-and-salty flavor.
- Reduce heat to a bare simmer. Pour beaten eggs into skillet in a thin, steady stream, holding chopsticks over edge of bowl to help distribute eggs evenly (see video above). Cover and cook until eggs are cooked to desired doneness, about 1 minute for runny eggs or 3 minutes for medium-firm.
- To Serve - transfer hot rice to a single large bowl or 2 individual serving bowls. Top with egg and chicken mixture, pouring out any excess broth from saucepan over rice. Add an extra egg yolk to center of each bowl, if desired (see note). Garnish with remaining sliced scallions and togarashi. Serve immediately.', '### Notes
- For a richer finished dish, use 7 eggs, reserving 2 of the yolks. Beat the extra egg whites together with the eggs in step 2, then add the reserved egg yolks to the finished bowls just before serving.
- consider adding mushrooms, experiment with this', '**Name:** Serious Eats
**URL:** https://www.seriouseats.com/oyakodon-japanese-chicken-and-egg-rice-bowl-recipe
**Type:** copy
**Modifications:** some, quantities and ingredients', 'croach', '2023-11-18 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('chicken');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '779b1147-1777-4b0a-b188-c60af0b2a2a7', id FROM tags WHERE name = 'chicken';
INSERT OR IGNORE INTO tags (name) VALUES ('egg');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '779b1147-1777-4b0a-b188-c60af0b2a2a7', id FROM tags WHERE name = 'egg';
INSERT OR IGNORE INTO tags (name) VALUES ('rice bowl');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '779b1147-1777-4b0a-b188-c60af0b2a2a7', id FROM tags WHERE name = 'rice bowl';

-- Recipe: Palak Paneer
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('6a79c089-5827-4c7a-9163-900ee13714e6', 'Palak Paneer', '', 'food', 'indian', '- 3 medium onions, diced
- 30g garlic, minced
- 30g ginger, minced
- 450g paneer
- 1.25kg frozen spinach
- 400g peas
- 1 green chilli, diced (deseed if bussy)
- 1 can chopped tomato
- 1.5kg chicken thigh
- 200g greek yogurt

### whole spices
- 0.5 tsp cumin seeds
- 0.5 tsp mustard seeds
- 12 cardamon pods
- 5 cloves
- 1 bay leaf

### spice mix
- 3 tsp salt
- 2 heaped tsp cumin
- 2 heaped tsp corriander
- 0.5 tsp fenugrek
- 1 tsp garam masala
- 0.25 tsp ground cinnamon
- 20 cracks black pepper
- pinch of MSG', '- In a larger pot (green ceramic), heat some oil.
- Brown the chicken in batches hitting it with some salt and pepper, just a few minutes each side. Remove chicken and put aside.
- Next to same pot of oil/chicken fat now add your whole spices and diced chilli and cook for a minute.
- Next add the onions, and cook em down good, around 10mins at least.
- Next add the ginger and garlic paste and cook for a couple minutes.
- Next add the powdered spice mix and cook for a minute. If at any point it''s looking dry add splashes of water. Then add the canned tomato and cook for another couple mins.
- Once the tamato juices have reduced a bit, add in the spinach. Cook it down.
- Blend the curry with an immersion blender or otherwise to your desired consistency.
- Add in the chicken, peas and paneer. Cover and cook until the chicken is done. Uncover if it''s too watery.
- Add yoghurt, mix and serve.', '### Notes
- Not all spices are absolutely mandatory. I''m really just chucking in whatever i have that smells good.', '**Name:** string
**URL:** string
**Type:** string
**Modifications:** string', 'string', '2025-08-14 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('curry');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '6a79c089-5827-4c7a-9163-900ee13714e6', id FROM tags WHERE name = 'curry';
INSERT OR IGNORE INTO tags (name) VALUES ('with rice');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '6a79c089-5827-4c7a-9163-900ee13714e6', id FROM tags WHERE name = 'with rice';

-- Recipe: Pork Dumplings
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('bcc2adc3-8514-4eda-beef-7deee6205c3a', 'Pork Dumplings', '', 'food', 'chinese', '- 3 * 500g packs dumpling wrappers (gyoza, wonton etc...)

### filling mixture
- 1kg pork mince, 5% fat
- 10 Dried Shiitake Mushrooms, rehydrated and finely chopped
- 15g Ginger, grated
- 2 cups Green Cabbage (half cabbage), finely chopped
- 8 green onions, sliced
- 1 can sliced Water Chestnuts (~133g drained), diced
- 2 eggs
- 1 tsp Salt
- 1/2 tsp Pepper
- 2 Tbsp Sesame Oil
- 2 Tbsp Soy Sauce
- 1 tsp Sugar

### cornstarch slurry
- 2 Tbsp Cornstarch
- 1/2 Cup Water', '- In a large bowl, combine all the filling mixture ingredients
- Create a cornstarch slurry and add it to the mixture and combine well
- Fill dumpling wrappers. Assuming you buy the usual gyoza/wonton wrappers, each dumpling wrapper is 12g. Add 15g of filling in the centre of the wrapper, then wet the edges and press them together to seal the mixture inside. Crimp the edges and repeat.
- To boil - boil for 5 to 8 minutes. Bit longer if frozen.
- To pan fry - pan fry on one side for a few minutes till a nice colour is developed. Then add some water and cover with a lid, steaming the tops.
- To freeze - place on trays or separated with glad wrap in the freezer so nothing sticks together. Then after an hour or two, combine them into desired size bags, they won''t stick together.', '### Notes
- savoy cabbage is a good choice
- Gyoza/wonton wrappers are generally 12g, so the 500g packets contain ~42 wrappers. By using 15g of mixture per wrapper, this recipe makes ~126 dumplings.

### Next
- could do with a teeny bit more salt', '**Name:** string
**URL:** https://www.maxiskitchen.com/blog/potstickers
**Type:** string
**Modifications:** string', 'string', '2024-09-29 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('pork');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'bcc2adc3-8514-4eda-beef-7deee6205c3a', id FROM tags WHERE name = 'pork';
INSERT OR IGNORE INTO tags (name) VALUES ('mince');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'bcc2adc3-8514-4eda-beef-7deee6205c3a', id FROM tags WHERE name = 'mince';
INSERT OR IGNORE INTO tags (name) VALUES ('freeze');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'bcc2adc3-8514-4eda-beef-7deee6205c3a', id FROM tags WHERE name = 'freeze';
INSERT OR IGNORE INTO tags (name) VALUES ('mealprep');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'bcc2adc3-8514-4eda-beef-7deee6205c3a', id FROM tags WHERE name = 'mealprep';
INSERT OR IGNORE INTO tags (name) VALUES ('dumpling');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'bcc2adc3-8514-4eda-beef-7deee6205c3a', id FROM tags WHERE name = 'dumpling';
INSERT OR IGNORE INTO tags (name) VALUES ('gyoza');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'bcc2adc3-8514-4eda-beef-7deee6205c3a', id FROM tags WHERE name = 'gyoza';
INSERT OR IGNORE INTO tags (name) VALUES ('wonton');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'bcc2adc3-8514-4eda-beef-7deee6205c3a', id FROM tags WHERE name = 'wonton';

-- Recipe: Pork Tenderloin Honey Garlic
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('cabb4578-6f7f-457e-9cb7-4b697853d538', 'Pork Tenderloin Honey Garlic', '', 'food', 'chinese', '- 2 pork tenderloin (pork fillet), 500g/1lb each (Note 1)
- 1 1/2 tbsp olive oil (or butter)
- 3 garlic cloves , very finely chopped

### Pork Tenderloin Rub
- 1 tsp garlic powder
- 1 tsp paprika
- 1/2 tsp salt
- 1/2 tsp black pepper

### Honey Garlic Sauce
- 3 tbsp cider vinegar (Note 2)
- 1 1/2 tbsp soy sauce, light or all purpose (Note 2)
- 1/2 cup honey (or maple syrup)', '- Preheat oven to 180C.
- Mix Sauce ingredients together.
- Mix Rub ingredients then sprinkle over the pork.
- Heat oil in a large oven proof skillet (Note 3) over high heat. Add pork and sear until golden all over.
- When pork is almost seared, push to the side, add garlic and cook until golden.
- Pour sauce in. Turn pork once, then immediately transfer to the oven.
- Roast 15 - 18 minutes or until the internal temperature is 68C / 155F (Note 4).
- Remove pork onto plate, cover loosely with foil and rest 5 minutes.
- Place skillet with sauce on stove over medium high heat, simmer rapidly for 3 minutes until liquid reduces down to thin syrup.
- Remove from stove, put pork in and turn to coat in sauce.
- Cut pork into thick slices and serve with sauce!', '', '**Name:** Recipe Tin Eats
**URL:** https://www.recipetineats.com/pork-tenderloin-with-honey-garlic-sauce/
**Type:** string
**Modifications:** string', 'string', '2024-06-08 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('asian');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'cabb4578-6f7f-457e-9cb7-4b697853d538', id FROM tags WHERE name = 'asian';
INSERT OR IGNORE INTO tags (name) VALUES ('garlic');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'cabb4578-6f7f-457e-9cb7-4b697853d538', id FROM tags WHERE name = 'garlic';
INSERT OR IGNORE INTO tags (name) VALUES ('honey');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'cabb4578-6f7f-457e-9cb7-4b697853d538', id FROM tags WHERE name = 'honey';
INSERT OR IGNORE INTO tags (name) VALUES ('pork');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'cabb4578-6f7f-457e-9cb7-4b697853d538', id FROM tags WHERE name = 'pork';

-- Recipe: Punjabi Chicken Curry
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('bc3ab227-8fc9-4f02-8b8d-58047774e606', 'Punjabi Chicken Curry', '', 'food', 'indian', '### Chicken', '- Combine the yogurt with the salt, spices, garlic, and ginger in a bowl, and mix well. Add the chicken thighs and make sure that the chicken is coated well in the yogurt marinade. Let it rest overnight, or 1 hour at least, in the fridge.
- Heat the oil in a pan. Add the cumin seeds, mustard seeds and bay leaf. Once they start to sizzle, add the onions and cook for 15 minutes, until a lovely golden color. Now add the tomatoes and cook for 10 minutes, until softened. Add the spices and salt, and cook for 1 minute. Add the marinated chicken, mix well, and cover. Cook for 40 to 45 minutes over low heat until the chicken is cooked through.
- Once you have cooked the curry, let it rest for 30 minutes to 1 hour. This makes the curry really intense and the chicken soaks up the flavors better. Sprinkle with some coriander and serve.', '', '**Name:** Food52
**URL:** https://food52.com/recipes/83814-punjabi-style-chicken-curry-recipe
**Type:** few changes
**Modifications:** string', 'dave', '2025-02-20 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('curry');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'bc3ab227-8fc9-4f02-8b8d-58047774e606', id FROM tags WHERE name = 'curry';
INSERT OR IGNORE INTO tags (name) VALUES ('chicken');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'bc3ab227-8fc9-4f02-8b8d-58047774e606', id FROM tags WHERE name = 'chicken';
INSERT OR IGNORE INTO tags (name) VALUES ('with rice');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'bc3ab227-8fc9-4f02-8b8d-58047774e606', id FROM tags WHERE name = 'with rice';

-- Recipe: Qeema Mince
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('c9b3da68-e2e9-4fd3-8a64-79578bb26699', 'Qeema Mince', '', 'food', 'indian', '- 3 tbsp vegetable or canola oil, or other neutral oil
- 18g ginger, finely minced or as paste
- 15g garlic, finely minced or as paste
- 1 large onion, finely diced
- 500g beef mince (ground beef)
- 1 tsp salt
- 1 tbsp kashmiri chilli
- 1 heaped tsp garam masala
- 1 heaped tsp cumin powder
- 1 heaped tsp coriander powder
- 1/2 tsp turmeric powder
- 1 cup water

### Garnish
- 1 green chilli (optional)
- coriander leaves (optional)', '- Heat oil in a skillet over high heat. Add onion and cook for 1 minute until it is starting to turn translucent. Add ginger and garlic and saute for 30 seconds until golden, don''t let it burn!
- Add beef and cook, breaking it up as you go, until it changes from pink to light brown. Add remaining ingredients EXCEPT water. Cook for a further 2 minutes to let the spices bloom.
- Add water, stir, then put a loose lid on and cook for 10 minutes. Turn heat down to medium and let it simmer for 10 minutes or until most of the water has evaporated.
- Serve over with basmati rice or plain white rice, garnish if desired.', '### Notes
- Most often now i make a double batch of this recipe at once, and also add 1 can of chickpeas and 1 can of blackbeans for some feg/fibre. I will sometimes add a bit more spices and salt to compensate for the extra vegetable mass.
- I often do a 50/50 beef/pork mix. Literally any mince will work here.', '**Name:** recipetineats
**URL:** https://www.recipetineats.com/qeema-indian-curried-beef/
**Type:** copy
**Modifications:** string', 'string', '2024-03-17 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('mince');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'c9b3da68-e2e9-4fd3-8a64-79578bb26699', id FROM tags WHERE name = 'mince';
INSERT OR IGNORE INTO tags (name) VALUES ('beef');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'c9b3da68-e2e9-4fd3-8a64-79578bb26699', id FROM tags WHERE name = 'beef';
INSERT OR IGNORE INTO tags (name) VALUES ('curry');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'c9b3da68-e2e9-4fd3-8a64-79578bb26699', id FROM tags WHERE name = 'curry';

-- Recipe: Sauce For Fried Rice
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('a71ab7bb-3e0b-4549-8f92-066cadbe6097', 'Sauce For Fried Rice', '', 'food', 'chinese', '- 1/3 cup Soy sauce
- 1/3 cup Oyster sauce
- 1/3 cup Mirin
- 1/4 cup Sesame oil
- 1 tablespoon white pepper', '- Combine all the ingredients in a jar. Keep sealed and store in the fridge.
- Use approximately 1-1/2 to 2 tablespoons of sauce for every cup of rice in your fried rice recipe. Taste and add more if desired.', '', '**Name:** Savor The Best
**URL:** https://savorthebest.com/sauce-for-fried-rice/
**Type:** copy
**Modifications:** minimal', 'croach', '2023-09-18 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('sauce');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'a71ab7bb-3e0b-4549-8f92-066cadbe6097', id FROM tags WHERE name = 'sauce';
INSERT OR IGNORE INTO tags (name) VALUES ('rice');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'a71ab7bb-3e0b-4549-8f92-066cadbe6097', id FROM tags WHERE name = 'rice';

-- Recipe: Smoked Salmon Pasta
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('07eba082-b037-4d15-b00c-9c8c7d828ee9', 'Smoked Salmon Pasta', '', 'food', 'western', '- 1 large brown onion, thinly sliced
- 2 tablespoons olive oil
- 500g button mushrooms, peeled and sliced
- 600 grams smoked salmon, roughly sliced
- 4-5 X 340ml Carnation Creamy Evaporated milk (75% less fat)
- ¾ of a container of Gourmet Garden Basil (from Woolworths)
- 1kg of Penne or preferred short pasta shape', '- Sautee onion in the olive oil with a little salt (to prevent the onion from burning) until soft
- Add the chopped mushrooms with a little more salt and pepper and cook until the liquid from the mushrooms evaporate
- Add the sliced smoked salmon and cook for a few minutes, stirring well (you will notice that the salmon changes colour and breaks up a little
- Add the basil into the sauce and stir through
- Add in the cream and stir the sauce often to prevent sticking to the bottom of the pan
- Cook, on medium heat, stirring often until the sauce thickens and coats the back of a wooden spoon
- Cook the pasta and add the sauce', '', '**Name:** Gina
**URL:** string
**Type:** string
**Modifications:** string', 'geenie', '2025-04-06 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('pasta');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '07eba082-b037-4d15-b00c-9c8c7d828ee9', id FROM tags WHERE name = 'pasta';
INSERT OR IGNORE INTO tags (name) VALUES ('smoked salmon');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '07eba082-b037-4d15-b00c-9c8c7d828ee9', id FROM tags WHERE name = 'smoked salmon';
INSERT OR IGNORE INTO tags (name) VALUES ('mushroom');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '07eba082-b037-4d15-b00c-9c8c7d828ee9', id FROM tags WHERE name = 'mushroom';
INSERT OR IGNORE INTO tags (name) VALUES ('creamy');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '07eba082-b037-4d15-b00c-9c8c7d828ee9', id FROM tags WHERE name = 'creamy';

-- Recipe: Sweet And Sour Fish
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('64057e6d-2862-4694-b44b-d441fdcc5b89', 'Sweet And Sour Fish', '', 'food', 'chinese', '### TODO: how to handle subheadings?

### For the fish:
- 900g fish cut into small portions (i used frozen basa)
- canola oil (for frying)
- 1.5 cup all-purpose flour
- 1/2 teaspoon baking powder
- 2 tablespoon cornstarch
- 1 teaspoon salt
- 1/2 teaspoon turmeric powder
- 1/4 teaspoon pepper (white, black, whatever)
- 1/2 teaspoon sesame oil
- 1 beer

### For the Sweet and Sour sauce:
- 1 red onion (cut into a 1-inch dice)
- 1 red capsicum (cut into a 1-inch dice)
- 1 green capsicum (cut into a 1-inch dice)
- 2 tablespoon ketchup
- 1.5 cup canned pineapple chunks
- 1.5 cup pineapple juice from the can
- 5 tablespoons red wine vinegar
- 2/3 water
- 1/2 teaspoon salt
- 4 tablespoons sugar
- 3 tablespoons cornstarch (mixed into a slurry with 4 tablespoons water)', '- Make sure your fish fillet is clean and pat dry to ensure your fried fish gets really crispy. Heat 3 cups of oil in a small pot (to conserve oil) to 380 degrees F. You can use a thermometer or check the temperature by putting a drop of batter into the oil. The batter should not turn brown right away. Instead, it should rise immediately to the surface and turn a very light golden brown.
- To make the batter, mix together all the dry ingredients - the flour, baking powder, cornstarch, salt, turmeric, and white pepper. When you''re ready to fry, mix in the sesame oil and cold seltzer water until the batter is smooth.
- Next, drop the fish fillets into the batter. Ensure they are evenly coated, but allow the excess to drip off (if there''s too much batter on the fish, you''ll end up with dough balls!). Carefully place the fish into the oil one piece at a time, ensuring that they don''t stick to each other. Fry in batches so that the fish pieces aren''t overcrowded, frying for 3-4 minutes or until golden brown.
- Scoop the fried fish out and transfer to a cooling rack placed over a baking sheet to drain. Repeat until all the fish has been fried.
- Next, put two teaspoons of the frying oil in a separate wok or skillet over high heat, and toss in the onions, and peppers. Stir-fry for 30 seconds, and then stir in the ketchup. Fry for another 20 seconds - frying the ketchup brings out the color and makes for a better depth of flavor in the sweet and sour sauce.
- Mix in the pineapple, pineapple juice, red wine vinegar, water, salt, and sugar, and bring the liquid to a low simmer for 2 minutes. With the sweet and sour sauce still simmering, slowly stir in the cornstarch slurry until the sauce is thick enough to coat a spoon.
- At this point, if the fried fish fillets have softened, you can refry them in the oil heated to 400F for 30 seconds in larger batches. The higher heat is needed because the oil will cool immediately after putting in a larger batch of fish in the oil.
- Once the fish is ready, toss it into the wok, and fold into the sauce with 3 or 4 scooping motions until the pieces are lightly coated. Plate and serve immediately!', '### Notes
- serve with rice

### Next
- try tripling the sauce ingredients
- do we need less capsicum? more onion?
- try adding some msg in the sauce', '**Name:** The Woks of Life
**URL:** https://thewoksoflife.com/sweet-sour-fish-fillet/
**Type:** copy
**Modifications:** some, quantities and ingredients', 'croach', '2023-11-18 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('pineapple');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '64057e6d-2862-4694-b44b-d441fdcc5b89', id FROM tags WHERE name = 'pineapple';
INSERT OR IGNORE INTO tags (name) VALUES ('fish');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '64057e6d-2862-4694-b44b-d441fdcc5b89', id FROM tags WHERE name = 'fish';
INSERT OR IGNORE INTO tags (name) VALUES ('fried');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '64057e6d-2862-4694-b44b-d441fdcc5b89', id FROM tags WHERE name = 'fried';
INSERT OR IGNORE INTO tags (name) VALUES ('frozen fish');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '64057e6d-2862-4694-b44b-d441fdcc5b89', id FROM tags WHERE name = 'frozen fish';

-- Recipe: Tomato Egg Prawn Stirfry
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('beef3c3d-24d2-42fc-9c55-7d9291e7f5c7', 'Tomato Egg Prawn Stirfry', '', 'food', 'chinese', '- 3 very large tomatoes
- 500g prawns
- 6 eggs
- 5 cloves of garlic, minced or chopped
- 4 spring onions, sliced and separated in two, half ingredient half garnish
- 2 tbsp oyster sauce
- 2 tsp sesame oil
- 2 tsp sugar
- 2 tsp chicken stock
- 1/2 tsp chicken salt
- 1/2 tsp pepper
- 2 tbsp corn starch
- 4 tbsp water', '- Put the rice cooker on in advance
- Mix together the sauce by combining oyster sauce, sesame oil, sugar, chicken stock, salt and pepper.
- Scramble the eggs, season lightly with salt and pepper
- In a pan with some oil, cook the eggs until just done, fluffy, moist. Take them off and put aside.
- Clean out the pan, then again with some oil, cook the garlic and half the spring onion for a minute.
- Add the prawns.
- When the prawns are almost done, add the sauce we made earlier and the tomatos. Cook until they have broken down and a release the juice.
- Make a slurry by combining the corn starch and water, and add it. Stir until thickened
- Add in the scrambled eggs.
- Plate ontop of rice, garnish with the remaining spring onion.
- Serve with rice', '### Notes
- We cooked this with prawns from frozen. The trick was we put the prawns BEFORE the galric and spring onion and let them cook on their own for a bit. Let them release water, drain it, then add more oil and the aromatics.

### Next
- try adding 1 or 2tbsp of tomato sauce for more tomato richness', '**Name:** Genius Eats
**URL:** https://www.therecipesource.com/the-recipe
**Type:** copy
**Modifications:** ratios, ingredients', 'croach', '2023-12-01 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('chinese');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'beef3c3d-24d2-42fc-9c55-7d9291e7f5c7', id FROM tags WHERE name = 'chinese';
INSERT OR IGNORE INTO tags (name) VALUES ('prawn');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'beef3c3d-24d2-42fc-9c55-7d9291e7f5c7', id FROM tags WHERE name = 'prawn';
INSERT OR IGNORE INTO tags (name) VALUES ('egg');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'beef3c3d-24d2-42fc-9c55-7d9291e7f5c7', id FROM tags WHERE name = 'egg';
INSERT OR IGNORE INTO tags (name) VALUES ('tomato');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'beef3c3d-24d2-42fc-9c55-7d9291e7f5c7', id FROM tags WHERE name = 'tomato';
INSERT OR IGNORE INTO tags (name) VALUES ('rice');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT 'beef3c3d-24d2-42fc-9c55-7d9291e7f5c7', id FROM tags WHERE name = 'rice';

-- Recipe: Unagi Sauce
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('623dde25-c032-451a-acff-2d20ad1c950e', 'Unagi Sauce', '', 'food', 'japanese', '- 60ml mirin
- 1½ Tbsp sake
- 2½ Tbsp sugar
- 60ml soy sauce', '- In a small saucepan, add the mirin, sake and sugar. Turn on the heat to medium and whisk all the ingredients together.
- Add soy sauce and bring to a boil. Once boiling, reduce the heat to low and continue simmering for 10 minutes. Toward the end of cooking, you will see more bubbles.
- Turn off the heat and let it cool. The sauce will thicken as it cools.', '### Notes
- A great way to dress up frozen unagi. Leftovers can be used to dress fried rice for a sweeter taste.', '**Name:** Just One Cookbook
**URL:** https://www.justonecookbook.com/homemade-unagi-sauce/
**Type:** copy
**Modifications:** minimal', 'croach', '2023-09-14 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('sauce');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '623dde25-c032-451a-acff-2d20ad1c950e', id FROM tags WHERE name = 'sauce';
INSERT OR IGNORE INTO tags (name) VALUES ('eel');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '623dde25-c032-451a-acff-2d20ad1c950e', id FROM tags WHERE name = 'eel';

-- Recipe: Vietnamese Noodle Salad
INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('550f0cd6-08c7-44df-ad85-7cf5424a6344', 'Vietnamese Noodle Salad', '', 'food', 'vietnamese', '### For the chicken & marinade:
- 450 g boneless, skinless chicken thighs
- 2 cloves garlic, minced
- 1 lime, juiced
- 30 ml fish sauce
- 15 ml soy sauce
- 25 g brown sugar
- 15 ml vegetable oil

### For the nuoc cham sauce:
- 3 cloves garlic, minced
- 1 lime, juiced
- 30 ml rice vinegar or white vinegar
- 60 ml fish sauce
- 38 g sugar
- 1 red chili, de-seeded and sliced, or substitute 2 teaspoons chili garlic sauce or Sriracha
- 120 ml cold water

### To assemble the bowls:
- 200g dried rice vermicelli noodles
- 250g bean sprouts
- 1 large carrot, julienned
- 1 seedless cucumber, julienned
- 2 baby cos lettuce, finely julienned
- 400g cherry tomatoes, halved
- Mint
- Corriander', '- In a medium bowl, combine the chicken thighs with your marinade ingredients, and set aside at room temperature for 30 mins to 1 hour while you prepare the other salad ingredients.
- Combine all the nuoc cham ingredients and stir until the sugar has completely dissolved into the sauce. Taste and adjust any of the ingredients if desired.
- Boil the rice vermicelli noodles according to the package instructions. Drain and rinse under cold running water. Set aside in a colander.
- Heat 2 tablespoons of vegetable oil in a cast iron skillet, grill pan, or frying pan over medium high heat. Sear the chicken for about 4 minutes per side, or until cooked through. Set aside on a plate.
- To assemble the salad, combine the rice noodles with bean sprouts, julienned carrots and cucumber, romaine lettuce, mint, and cilantro. Slice the chicken thighs and add to the salad. Serve with your nuoc cham sauce.', '### Notes
- Also excellent without any meat as a side dish
- Would work great with any meat, but particularly beef

### Next
- We haven''t actually tried this with the chicken yet, do that!', '**Name:** The Woks of Life
**URL:** https://thewoksoflife.com/vietnamese-rice-noodle-salad-chicken/
**Type:** string
**Modifications:** removed meat', 'meggles', '2024-03-24 00:00:00', '2025-11-04 18:03:24');
INSERT OR IGNORE INTO tags (name) VALUES ('salad');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '550f0cd6-08c7-44df-ad85-7cf5424a6344', id FROM tags WHERE name = 'salad';
INSERT OR IGNORE INTO tags (name) VALUES ('noodle');
INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '550f0cd6-08c7-44df-ad85-7cf5424a6344', id FROM tags WHERE name = 'noodle';

COMMIT;

-- migrate:up
CREATE TABLE categories (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    icon       TEXT NOT NULL,
    color      TEXT NOT NULL,
    budget     DECIMAL(10, 2) NOT NULL DEFAULT 0,
    keywords   TEXT[] NOT NULL DEFAULT '{}',
    sort_order INT NOT NULL DEFAULT 0
);

INSERT INTO categories (name, icon, color, budget, keywords, sort_order) VALUES
    ('Groceries',     'pi pi-shopping-cart', '#f5c518', 400, ARRAY['grocer','supermarket','food','market'], 1),
    ('Dining',        'pi pi-utensils',      '#f97316', 300, ARRAY['restaurant','cafe','coffee','dining','eat','lunch','dinner','breakfast'], 2),
    ('Transport',     'pi pi-car',           '#3b82f6', 200, ARRAY['uber','lyft','bus','train','transit','fuel','gas','transport','taxi'], 3),
    ('Bills',         'pi pi-file-edit',     '#a855f7', 500, ARRAY['bill','electric','water','internet','phone','rent','insurance'], 4),
    ('Shopping',      'pi pi-tag',           '#ec4899', 250, ARRAY['shop','amazon','clothing','clothes','shoes','mall'], 5),
    ('Entertainment', 'pi pi-star',          '#14b8a6', 150, ARRAY['netflix','spotify','movie','game','entertain','stream'], 6),
    ('Health',        'pi pi-heart',         '#ef4444', 100, ARRAY['gym','health','pharmacy','doctor','medical','fitness'], 7),
    ('Other',         'pi pi-ellipsis-h',    '#71717a', 200, ARRAY[]::TEXT[], 8);

-- migrate:down
DROP TABLE IF EXISTS categories;
